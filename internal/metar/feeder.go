package metar

import (
	"github.com/rs/zerolog/log"
	"github.com/skybi/nuntius/internal/client"
	"github.com/skybi/nuntius/internal/queue"
	"time"
)

// Feeder represents the worker queueing and feeding new METARs assembled by the cycle workers
type Feeder struct {
	apiClient *client.Client
	batchSize int
	interval  time.Duration

	queue    *queue.Queue[string]
	priority []string

	running bool
	stop    chan struct{}
}

// NewFeeder creates a new METAR feeder
func NewFeeder(apiClient *client.Client, batchSize int, interval time.Duration) *Feeder {
	return &Feeder{
		apiClient: apiClient,
		batchSize: batchSize,
		interval:  interval,
		queue:     queue.New[string](),
	}
}

// Queue queues METARs to feed
func (feeder *Feeder) Queue(metars []string) {
	feeder.queue.Push(metars...)
}

// Start starts the feeding task
func (feeder *Feeder) Start() {
	if feeder.running {
		return
	}
	feeder.running = true
	feeder.stop = make(chan struct{})
	go func() {
		for {
			select {
			case <-feeder.stop:
				return
			case <-time.After(feeder.interval):
				if len(feeder.priority) > 0 {
					values := feeder.priority[:feeder.batchSize]
					_, err := feeder.apiClient.FeedMETARs(values)
					if err != nil {
						log.Err(err).Msg("could not feed prioritized METARs")
						continue
					}
					log.Debug().Int("amount", len(values)).Msg("fed prioritized METARs")
					feeder.priority = feeder.priority[feeder.batchSize:]
				} else if feeder.queue.Size() > 0 {
					values := feeder.queue.PopN(feeder.batchSize)
					_, err := feeder.apiClient.FeedMETARs(values)
					if err != nil {
						feeder.priority = values
						log.Err(err).Msg("could not feed METARs; prioritizing them for the next run")
						continue
					}
					log.Debug().Int("amount", len(values)).Msg("fed METARs")
				}
			}
		}
	}()
}

// Stop stops the feeding task
func (feeder *Feeder) Stop() {
	if !feeder.running {
		return
	}
	close(feeder.stop)
	feeder.running = false
}
