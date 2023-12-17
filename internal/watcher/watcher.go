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

func (w LogWatcher) shouldSendMessage(logLine notifier.LogLine) bool {
	switch {
	case logLine.EventType == notifier.LoggedIn && w.WatchSettings.WatchAcceptedLogins:
		return true
	case logLine.EventType == notifier.FailedLoginAttempt && w.WatchSettings.WatchFailedLogins:
		return true
	case logLine.EventType == notifier.FailedLoginAttemptInvalidUsername && w.WatchSettings.WatchFailedLoginInvalidUsername:
		return true
	default:
		return false
	}
}

// TODO(mgottlieb) refactor this into more unit-testable funcs
func (w LogWatcher) Watch() {
	currentSize := 0
	for {
		file, err := os.Open(w.LogFile)
		if err != nil {
			log.Fatal().Err(err)
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			log.Fatal().Err(err)
		}

		if stat.Size() > int64(currentSize) {
			var currentLine int

			if _, err := os.Stat(stateFile); err == nil {
				state, err := os.Open(stateFile)
				if err != nil {
					log.Fatal().Err(err)
				}
				defer state.Close()

				scanner := bufio.NewScanner(state)
				for scanner.Scan() {
					currentLine, err = strconv.Atoi(scanner.Text())
					if err != nil {
						log.Fatal().Err(err)
					}
				}

			}

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)
			for i := 0; scanner.Scan(); i++ {
				if i <= currentLine {
					continue
				}
				line := scanner.Text()
				logLine := ParseLogLine(line)

				if w.shouldSendMessage(logLine) {
					if err := w.Notifier.Notify(logLine); err != nil {
						log.Error().Err(err)
						continue
					} else {
						log.Info().Msg("Posted message to slack")
					}
				}

				state, err := os.Create(stateFile)
				if err != nil {
					log.Error().Err(err)
					continue
				}
				if _, err := state.WriteString(fmt.Sprintf("%d", i)); err != nil {
					log.Error().Err(err)
					continue
				}

				if err := state.Sync(); err != nil {
					log.Error().Err(err)
					continue
				}

			}
			currentSize = int(stat.Size())
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
		case strings.Contains(line, "Failed password"), strings.Contains(line, "Connection closed by authenticating user"):
			logLine.EventType = notifier.FailedLoginAttempt
		case strings.Contains(strings.ToLower(line), "invalid user"):
			logLine.EventType = notifier.FailedLoginAttemptInvalidUsername
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
