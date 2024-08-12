package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/hpcloud/tail"
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

// reader is the interface that wraps the basic Read method.
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . reader
type reader interface {
	Read(p []byte) (n int, err error)
}

// file is the interface for interacting with the filesystem.
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . file
type file interface {
	Open(name string) (*os.File, error)
}

func New(logFile string, notifier notifierClient, hostMachine string, watchSettings config.WatchSettings, processedLineTracker processedLineTracker, file file) App {
	return App{
		logFile:              logFile,
		notifier:             notifier,
		hostMachine:          hostMachine,
		watchSettings:        watchSettings,
		processedLineTracker: processedLineTracker,
		file:                 file,
	}
}

type App struct {
	logFile              string
	notifier             notifierClient
	hostMachine          string
	watchSettings        config.WatchSettings
	processedLineTracker processedLineTracker
	file                 file
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

func (a App) processLine(line string, lineNumber int) error {
	logLine := a.parseLogLine(line)
	if !a.shouldSendMessage(logLine.EventType) {
		return nil
	}

	if err := a.notifier.Notify(logLine); err != nil {
		return fmt.Errorf("error sending notification: %w", err)
	}

	log.Info().Msg("notification message sent")
	err := a.processedLineTracker.UpdateLastProcessedLine(lineNumber)
	if err != nil {
		log.Error().Err(err)
		return fmt.Errorf("%v: %w", "failed updating last processed line", err)
	}
	return nil
}

func (a App) Watch() error {
	t, err := tail.TailFile(a.logFile, tail.Config{Follow: true})
	if err != nil {
		return fmt.Errorf("error tailing file: %w", err)
	}
	for line := range t.Lines {
		if err := a.processLine(line.Text, 0); err != nil {
			return fmt.Errorf("error processing line: %w", err)
		}
	}
	return nil
}
