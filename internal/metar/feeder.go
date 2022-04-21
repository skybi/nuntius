package metar

import (
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/skybi/nuntius/internal/client"
	"github.com/skybi/nuntius/internal/file"
	"github.com/skybi/nuntius/internal/queue"
	"os"
	"time"
)

var queueBackupFilepath = "./data/metar/feeder-queue"

// Feeder represents the worker queueing and feeding new METARs assembled by the cycle workers
type Feeder struct {
	queue *queue.Queue[string]

	apiClient *client.Client
	batchSize int

	interval time.Duration
	running  bool
	stop     chan struct{}
}

// NewFeeder creates a new METAR feeder
func NewFeeder(apiClient *client.Client, batchSize int, interval time.Duration) *Feeder {
	return &Feeder{
		queue:     queue.New[string](),
		apiClient: apiClient,
		batchSize: batchSize,
		interval:  interval,
	}
}

// Queue queues METARs to feed and fixes them beforehand
func (feeder *Feeder) Queue(metars []string) {
	for i, metar := range metars {
		metars[i] = fix(metar)
	}
	feeder.queue.Push(metars...)
}

// Start starts the feeding task
func (feeder *Feeder) Start() error {
	if feeder.running {
		return nil
	}

	if err := feeder.restoreQueue(); err != nil {
		return err
	}

	feeder.running = true
	feeder.stop = make(chan struct{})
	go func() {
		for {
			select {
			case <-feeder.stop:
				return
			case <-time.After(feeder.interval):
				if feeder.queue.Size() > 0 {
					values := feeder.queue.PopN(feeder.batchSize)
					err := feeder.apiClient.FeedMETARsRelaxed(values)
					if err != nil {
						feeder.queue.Push(values...)
						log.Err(err).Msg("could not feed METARs; appending them to the queue again")
						continue
					}
					log.Debug().Int("amount", len(values)).Msg("fed METARs")
				}
			}
		}
	}()
	return nil
}

// Stop stops the feeding task
func (feeder *Feeder) Stop() error {
	if !feeder.running {
		return nil
	}
	close(feeder.stop)
	feeder.running = false

	return feeder.backupQueue()
}

func (feeder *Feeder) backupQueue() error {
	data, err := queue.Serialize(feeder.queue)
	if err != nil {
		return err
	}

	return file.Write(queueBackupFilepath, data)
}

func (feeder *Feeder) restoreQueue() error {
	data, err := file.Read(queueBackupFilepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	restored, err := queue.Deserialize[string](data)
	if err != nil {
		return err
	}
	feeder.queue = restored
	return nil
}
