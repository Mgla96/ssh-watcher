package linetracker

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

func NewFileProcessedLineTracker(stateFilePath string) FileProcessedLineTracker {
	return FileProcessedLineTracker{
		StateFilePath: stateFilePath,
	}
}

type FileProcessedLineTracker struct {
	StateFilePath string
}

// GetLastProcessedLine reads the statefile and extracts the last processed line number
// in the ssh log file.
func (f FileProcessedLineTracker) GetLastProcessedLine() (int, error) {
	if _, err := os.Stat(f.StateFilePath); os.IsNotExist(err) {
		return 0, err
	}

	state, err := os.Open(f.StateFilePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed opening state file")
		return 0, err
	}
	defer state.Close()

	scanner := bufio.NewScanner(state)
	var currentLine int
	for scanner.Scan() {
		currentLine, err = strconv.Atoi(scanner.Text())
		if err != nil {
			log.Error().Err(err).Msg("Failed converting state file line to int")
			return 0, err
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Msg("Error while reading state file")
		return 0, err
	}

	return currentLine, nil
}

func (f FileProcessedLineTracker) UpdateLastProcessedLine(lineNumber int) error {
	state, err := os.Create(f.StateFilePath)
	if err != nil {
		message := "failed to create or truncate state file"
		log.Error().Err(err).Msg(message)
		return fmt.Errorf("%v: %w", message, err)
	}
	defer state.Close()

	if _, err := state.WriteString(fmt.Sprintf("%d", lineNumber)); err != nil {
		message := "failed to write to state file"
		log.Error().Err(err).Msg(message)
		return fmt.Errorf("%v: %w", message, err)
	}

	if err := state.Sync(); err != nil {
		message := "failed to syc state file"
		log.Error().Err(err).Msg(message)
		return fmt.Errorf("%v: %w", message, err)
	}

	return nil
}
