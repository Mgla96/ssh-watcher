package watcher

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/mgla96/ssh-watcher/internal/notifier"

	"github.com/rs/zerolog/log"
)

type WatchSettings struct {
	WatchAcceptedLogins             bool
	WatchFailedLogins               bool
	WatchFailedLoginInvalidUsername bool
	WatchSleepInterval              time.Duration
}

func NewLogWatcher(logFile string, notifier notifier.Notifier, hostMachine string, watchSettings WatchSettings, processedLineTracker ProcessedLineTracker) LogWatcher {
	return LogWatcher{
		LogFile:              logFile,
		Notifier:             notifier,
		HostMachine:          hostMachine,
		WatchSettings:        watchSettings,
		ProcessedLineTracker: processedLineTracker,
	}
}

type ProcessedLineTracker interface {
	GetLastProcessedLine() (int, error)
	UpdateLastProcessedLine(lineNumber int) error
}

type LogWatcher struct {
	LogFile              string
	Notifier             notifier.Notifier
	HostMachine          string
	WatchSettings        WatchSettings
	ProcessedLineTracker ProcessedLineTracker
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

func (w LogWatcher) parseLogLine(line string) notifier.LogLine {
	logLine := notifier.LogLine{}
	if !strings.Contains(line, "sshd") {
		return logLine
	}

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

	return logLine
}

func (w LogWatcher) processNewLogLines(file *os.File, lastProcessedLine int) error {
	scanner := bufio.NewScanner(file)
	for lineNumber := 0; scanner.Scan(); lineNumber++ {
		// TODO(mgottlieb) we do not need to scan from very beginning line every time.
		if lineNumber <= lastProcessedLine {
			continue
		}

		line := scanner.Text()
		logLine := w.parseLogLine(line)

		if w.shouldSendMessage(logLine.EventType) {
			if err := w.Notifier.Notify(logLine); err != nil {
				log.Error().Err(err)
				continue
			} else {
				log.Info().Msg("notification message sent")
			}
		}

		err := w.ProcessedLineTracker.UpdateLastProcessedLine(lineNumber)
		if err != nil {
			log.Error().Err(err)
			return fmt.Errorf("%v: %w", "failed updating last processed line", err)
		}
	}
	return nil
}

// TODO(mgottlieb) refactor this into more unit-testable funcs.
func (w LogWatcher) Watch() error {
	file, err := os.Open(w.LogFile)
	if err != nil {
		return fmt.Errorf("error opening log file: %w", err)
	}
	defer file.Close()

	var lastProcessedOffset int64 = 0
	lastFileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error returning last file info: %w", err)
	}

	for {
		stat, err := file.Stat()
		if err != nil {
			return fmt.Errorf("error returning file info: %w", err)
		}
		// TODO(mgottlieb) check log rotation.
		if isLogRotated(stat, lastFileInfo) {
			if err := file.Close(); err != nil {
				return fmt.Errorf("error closing file: %w", err)
			}
			file, err = os.Open(w.LogFile)
			if err != nil {
				return fmt.Errorf("error opening file: %w", err)
			}
			stat, err = file.Stat()
			if err != nil {
				return fmt.Errorf("error returning file info when log rotated: %w", err)
			}
			lastProcessedOffset = 0
			lastFileInfo = stat
		}

		if stat.Size() > lastProcessedOffset {
			lastProcessedLine, err := w.ProcessedLineTracker.GetLastProcessedLine()
			if err != nil {
				return fmt.Errorf("error getting last processed line: %w", err)
			}

			err = w.processNewLogLines(file, lastProcessedLine)
			if err != nil {
				return fmt.Errorf("error processing new log lines: %w", err)
			}
			lastProcessedOffset = stat.Size()
		}

		time.Sleep(w.WatchSettings.WatchSleepInterval)
	}
}

func isLogRotated(currentFileInfo fs.FileInfo, lastFileInfo fs.FileInfo) bool {
	return !os.SameFile(currentFileInfo, lastFileInfo)
}
