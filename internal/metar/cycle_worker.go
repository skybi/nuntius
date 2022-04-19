package metar

import (
	"bufio"
	"errors"
	"github.com/jlaffaye/ftp"
	"github.com/rs/zerolog/log"
	"github.com/skybi/nuntius/internal/file"
	"github.com/skybi/nuntius/internal/set"
	"io"
	"os"
	"time"
)

const (
	ftpAddress        = "tgftp.nws.noaa.gov:21"
	ftpCyclesLocation = "/data/observations/metar/cycles/"
)

// cycleWorker represents a worker fetching, deduplicating and queuing a single METAR cycle of the NOAA's FTP data server
type cycleWorker struct {
	ftpConn        *ftp.ServerConn
	remoteFileName string
	lastChanged    time.Time

	feeder        *Feeder
	stateFilePath string

	running  bool
	stopChan chan struct{}
}

func (worker *cycleWorker) start() error {
	if worker.running {
		return nil
	}

	ftpConn, err := openFTPConn()
	if err != nil {
		return err
	}
	worker.ftpConn = ftpConn

	worker.running = true
	worker.stopChan = make(chan struct{})
	go func() {
		for {
			select {
			case <-worker.stopChan:
				return
			case <-time.After(30 * time.Second):
				// We only want to process files that were modified since we processed them the previous time
				var lastChanged time.Time
				if worker.ftpConn.IsGetTimeSupported() {
					changed, err := worker.ftpConn.GetTime(worker.remoteFileName)
					if err != nil {
						log.Error().Err(err).Msg("could not check modification time")
						continue
					}
					if changed == worker.lastChanged {
						continue
					}
					lastChanged = changed
				}

				// Open a data connection to read the file
				reader, err := worker.ftpConn.Retr(worker.remoteFileName)
				if err != nil {
					log.Error().Err(err).Msg("could not read remote file")
					continue
				}

				// Extract and deduplicate the raw METARs out of the file
				metars, err := extractMETARs(bufio.NewReader(reader))
				if err != nil {
					reader.Close()
					log.Error().Err(err).Msg("could not extract METARs out of remote file")
					continue
				}
				reader.Close()

				// Load the METARs we processed the previous time
				state, err := loadCycleState(worker.stateFilePath)
				if err != nil {
					log.Error().Err(err).Msg("could not read current METAR state")
					continue
				}

				// Update the state
				if err := saveCycleState(worker.stateFilePath, metars); err != nil {
					log.Error().Err(err).Msg("could not update current METAR state")
					continue
				}

				// Add the difference between both sets to the feeding queue
				values := set.Diff(metars, state).ToSlice()
				worker.feeder.Queue(values)
				log.Debug().Int("amount", len(values)).Msg("queued METARs to feed")

				worker.lastChanged = lastChanged
			}
		}
	}()

	return nil
}

func (worker *cycleWorker) stop() {
	if !worker.running {
		return
	}
	close(worker.stopChan)
	worker.ftpConn.Quit()
	worker.running = false
}

func openFTPConn() (*ftp.ServerConn, error) {
	ftpConn, err := ftp.Dial(ftpAddress, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}
	if err := ftpConn.Login("anonymous", "anonymous"); err != nil {
		ftpConn.Quit()
		return nil, err
	}
	if err := ftpConn.ChangeDir(ftpCyclesLocation); err != nil {
		ftpConn.Quit()
		return nil, err
	}
	return ftpConn, nil
}

func extractMETARs(reader *bufio.Reader) (*set.HashSet[string], error) {
	hashSet := set.NewHashSet[string]()
	if reader.Size() == 0 {
		return hashSet, nil
	}

	end := false
	for !end {
		line, isPrefix, err := reader.ReadLine()
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}
		for isPrefix {
			fragment, isAnotherPrefix, err := reader.ReadLine()
			if err != nil && !errors.Is(err, io.EOF) {
				return nil, err
			}
			line = append(line, fragment...)
			isPrefix = isAnotherPrefix
		}
		end = err != nil && errors.Is(err, io.EOF)

		// 47 = '/'
		if len(line) == 0 || (len(line) > 4 && line[4] == 47) {
			continue
		}

		hashSet.Add(string(line))
	}

	return hashSet, nil
}

func loadCycleState(filepath string) (*set.HashSet[string], error) {
	data, err := file.Read(filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return set.NewHashSet[string](), nil
		}
		return nil, err
	}

	metars, err := set.Deserialize[string](data)
	if err != nil {
		return nil, err
	}
	return metars, nil
}

func saveCycleState(filepath string, metars *set.HashSet[string]) error {
	data, err := set.Serialize(metars)
	if err != nil {
		return err
	}

	return file.Write(filepath, data)
}
