package app

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/mgla96/ssh-watcher/internal/config"
	"github.com/mgla96/ssh-watcher/internal/notifier"

	"github.com/rs/zerolog/log"
)

// notifierClient is an interface for sending notifications
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . notifierClient
type notifierClient interface {
	Notify(LogLine notifier.LogLine) error
}

// processedLineTracker is an interface for tracking the last processed line
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . processedLineTracker
type processedLineTracker interface {
	GetLastProcessedLine() (int, error)
	UpdateLastProcessedLine(lineNumber int) error
}

func NewApp(logFile string, notifier notifierClient, hostMachine string, watchSettings config.WatchSettings, processedLineTracker processedLineTracker) App {
	return App{
		logFile:              logFile,
		notifier:             notifier,
		hostMachine:          hostMachine,
		watchSettings:        watchSettings,
		processedLineTracker: processedLineTracker,
	}
}

type App struct {
	logFile              string
	notifier             notifierClient
	hostMachine          string
	watchSettings        config.WatchSettings
	processedLineTracker processedLineTracker
}

func (a App) shouldSendMessage(eventType notifier.EventType) bool {
	switch {
	case eventType == notifier.LoggedIn && a.watchSettings.AcceptedLogins:
		return true
	case eventType == notifier.FailedLoginAttempt && a.watchSettings.FailedLogins:
		return true
	case eventType == notifier.FailedLoginAttemptInvalidUsername && a.watchSettings.FailedLoginInvalidUsername:
		return true
	default:
		return false
	}
}

func (a App) parseLogLine(line string) notifier.LogLine {
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

func (a App) processNewLogLines(file *os.File, lastProcessedLine int) error {
	scanner := bufio.NewScanner(file)
	for lineNumber := 0; scanner.Scan(); lineNumber++ {
		// TODO(mgottlieb) we do not need to scan from very beginning line every time.
		if lineNumber <= lastProcessedLine {
			continue
		}

		line := scanner.Text()
		logLine := a.parseLogLine(line)

		if a.shouldSendMessage(logLine.EventType) {
			if err := a.notifier.Notify(logLine); err != nil {
				log.Error().Err(err)
				continue
			} else {
				log.Info().Msg("notification message sent")
			}
		}

		err := a.processedLineTracker.UpdateLastProcessedLine(lineNumber)
		if err != nil {
			log.Error().Err(err)
			return fmt.Errorf("%v: %w", "failed updating last processed line", err)
		}
	}
	return nil
}

// TODO(mgottlieb) refactor this into more unit-testable funcs.
func (a App) Watch() error {
	file, err := os.Open(a.logFile)
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
			file, err = os.Open(a.logFile)
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
			lastProcessedLine, err := a.processedLineTracker.GetLastProcessedLine()
			if err != nil {
				return fmt.Errorf("error getting last processed line: %w", err)
			}

			err = a.processNewLogLines(file, lastProcessedLine)
			if err != nil {
				return fmt.Errorf("error processing new log lines: %w", err)
			}
			lastProcessedOffset = stat.Size()
		}

		time.Sleep(time.Duration(a.watchSettings.SleepInterval) * time.Second)
	}
}

func isLogRotated(currentFileInfo fs.FileInfo, lastFileInfo fs.FileInfo) bool {
	return !os.SameFile(currentFileInfo, lastFileInfo)
}
