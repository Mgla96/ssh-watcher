package watcher

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mgla96/ssh-watcher/internal/notifier"

	"github.com/rs/zerolog/log"
)

const (
	// stateFile keeps track of the last processed line by ssh watcher so restarts
	// of the service do not reprocess all ssh history.
	stateFile = "/tmp/authlog-state"
)

type WatchSettings struct {
	WatchAcceptedLogins             bool
	WatchFailedLogins               bool
	WatchFailedLoginInvalidUsername bool
}

type LogWatcher struct {
	LogFile       string
	Notifier      notifier.Notifier
	HostMachine   string
	WatchSettings WatchSettings
}

func (w LogWatcher) shouldSendMessage(eventType notifier.EventType) bool {
	switch {
	case eventType == notifier.LoggedIn && w.WatchSettings.WatchAcceptedLogins:
		return true
	case eventType == notifier.FailedLoginAttempt && w.WatchSettings.WatchFailedLogins:
		return true
	case eventType == notifier.FailedLoginAttemptInvalidUsername && w.WatchSettings.WatchFailedLoginInvalidUsername:
		return true
	default:
		return false
	}
}

// getLastProcessedLine reads the statefile and extracts the last processed line number
// in the ssh log file.
func (w LogWatcher) getLastProcessedLine() int {
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		return 0
	}

	state, err := os.Open(stateFile)
	if err != nil {
		log.Error().Err(err).Msg("Failed opening state file")
		return 0
	}
	defer state.Close()

	scanner := bufio.NewScanner(state)
	var currentLine int
	for scanner.Scan() {
		currentLine, err = strconv.Atoi(scanner.Text())
		if err != nil {
			log.Error().Err(err).Msg("Failed converting state file line to int")
			return 0
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Msg("Error while reading state file")
		return 0
	}

	return currentLine

}

func (w LogWatcher) updateLastProcessedLine(lineNumber int) error {
	state, err := os.Create(stateFile)
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

// TODO(mgottlieb) refactor this into more unit-testable funcs.
func (w LogWatcher) Watch() {
	file, err := os.Open(w.LogFile)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer file.Close()

	var lastProcessedOffset int64 = 0

	for {
		stat, err := file.Stat()
		if err != nil {
			log.Fatal().Err(err)
		}
		// TODO(mgottlieb) check log rotation.
		// if stat.Size() < lastProcessedOffset || stat.ModTime().After(lastFileInfo.ModTime()) {}

		if stat.Size() > lastProcessedOffset {
			lastProcessedLine := w.getLastProcessedLine()
			scanner := bufio.NewScanner(file)
			for lineNumber := 0; scanner.Scan(); lineNumber++ {
				// TODO(mgottlieb) we do not need to scan from very beginning line every time.
				if lineNumber <= lastProcessedLine {
					continue
				}

				line := scanner.Text()
				logLine := ParseLogLine(line)

				if w.shouldSendMessage(logLine.EventType) {
					if err := w.Notifier.Notify(logLine); err != nil {
						log.Error().Err(err)
						continue
					} else {
						log.Info().Msg("Posted message to slack")
					}
				}

				err := w.updateLastProcessedLine(lineNumber)
				if err != nil {
					log.Error().Err(err)
					continue
				}
			}
			lastProcessedOffset = stat.Size()
		}

		time.Sleep(2 * time.Second)
	}
}

func ParseLogLine(line string) notifier.LogLine {
	logLine := notifier.LogLine{}
	if strings.Contains(line, "sshd") {
		switch {
		case strings.Contains(line, "Accepted password"), strings.Contains(line, "Accepted publickey"):
			logLine.EventType = notifier.LoggedIn
		case strings.Contains(strings.ToLower(line), "invalid user"):
			logLine.EventType = notifier.FailedLoginAttemptInvalidUsername
		case strings.Contains(line, "Failed password"), strings.Contains(line, "Connection closed by authenticating user"):
			logLine.EventType = notifier.FailedLoginAttempt
		}

		if logLine.EventType != "" {
			parts := strings.Split(line, " ")
			logLine.LoginTime = parts[0] + " " + parts[1]
			for i, part := range parts {
				if part == "from" {
					logLine.IpAddress = parts[i+1]
				}
				if part == "user" || part == "for" {
					logLine.Username = parts[i+1]
				}
			}
		}
	}
	return logLine
}
