package linetracker

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog/log"
)

func NewFileProcessedLineTracker(stateFilePath string) FileProcessedLineTracker {
	return FileProcessedLineTracker{
		StateFilePath: stateFilePath,
	}
}

func handleError(message string, err error) error {
	log.Error().Err(err).Msg(message)
	return fmt.Errorf("%v: %w", message, err)
}

type FileProcessedLineTracker struct {
	StateFilePath string
}

// GetLastProcessedLine reads the statefile and extracts the last processed line number
// in the ssh log file.
func (f FileProcessedLineTracker) GetLastProcessedLine() (int, error) {
	dir := filepath.Dir(f.StateFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return 0, handleError("failed to create directory for state file", err)
	}

	state, err := os.OpenFile(f.StateFilePath, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return 0, handleError("failed opening or creating state file", err)
	}
	defer state.Close()

	scanner := bufio.NewScanner(state)
	var currentLine int
	for scanner.Scan() {
		currentLine, err = strconv.Atoi(scanner.Text())
		if err != nil {
			return 0, handleError("failed converting state file line to int", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("%v: %w", "error while reading state file", err)
	}

	return currentLine, nil
}

func (f FileProcessedLineTracker) UpdateLastProcessedLine(lineNumber int) error {
	state, err := os.Create(f.StateFilePath)
	if err != nil {
		return handleError("failed to create or truncate state file", err)
	}
	defer state.Close()

	if _, err := state.WriteString(fmt.Sprintf("%d", lineNumber)); err != nil {
		return handleError("failed to write to state file", err)
	}

	if err := state.Sync(); err != nil {
		return handleError("failed to syc state file", err)
	}

	return nil
}
