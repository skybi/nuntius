package metar

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"path/filepath"
)

// CycleWorkers groups and controls the 24 needed cycle workers instances
type CycleWorkers struct {
	workers [24]*cycleWorker
	running bool
}

// InitWorkers creates, initializes and groups the 24 different METAR cycle workers
func InitWorkers(feeder *Feeder) *CycleWorkers {
	// Create the 24 workers
	var workers [24]*cycleWorker
	for i := 0; i < 24; i++ {
		path, _ := filepath.Abs(fmt.Sprintf("./data/metar/cycle-state-%02d", i))
		workers[i] = &cycleWorker{
			remoteFileName: fmt.Sprintf("%02dZ.TXT", i),
			feeder:         feeder,
			stateFilePath:  path,
		}
	}

	return &CycleWorkers{
		workers: workers,
		running: false,
	}
}

// Start starts all cycle workers
func (workers *CycleWorkers) Start() {
	if workers.running {
		return
	}
	for i, worker := range workers.workers {
		if err := worker.start(); err != nil {
			log.Error().Err(err).Int("worker", i).Msg("could not start worker")
		}
	}
	workers.running = true
}

// Stop stops all cycle workers
func (workers *CycleWorkers) Stop() {
	if !workers.running {
		return
	}
	for _, worker := range workers.workers {
		worker.start()
	}
	workers.running = false
}
