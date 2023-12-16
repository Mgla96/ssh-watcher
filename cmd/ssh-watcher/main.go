package main

import (
	"fmt"
	"os"

	"github.com/mgla96/ssh-watcher/internal/notifier"
	"github.com/mgla96/ssh-watcher/internal/watcher"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Starting watcher")

	hostUrl := os.Getenv("HOST_URL")
	webhookUrl := os.Getenv("SLACK_WEBHOOK_URL")
	notifier := notifier.SlackNotifier{
		WebhookURL:    webhookUrl,
		SlackChannel:  "#ssh-alerts",
		SlackUsername: "poe-ssh-bot",
		SlackIcon:     ":ghost:",
	}

	logFileLocation := os.Getenv("WATCH_LOGFILE")
	if len(logFileLocation) == 0 {
		logFileLocation = "/var/log/auth.log"
	}
	log.Info().Msg(fmt.Sprintf("webhook url: %s, logfile: %s", webhookUrl, logFileLocation))
	watcher := watcher.LogWatcher{
		LogFile:     logFileLocation,
		Notifier:    notifier,
		HostMachine: hostUrl,
	}
	watcher.Watch()
}
