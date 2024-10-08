package main

import (
	"fmt"

	"github.com/mgla96/ssh-watcher/internal/app"
	"github.com/mgla96/ssh-watcher/internal/config"
	"github.com/mgla96/ssh-watcher/internal/file"
	"github.com/mgla96/ssh-watcher/internal/linetracker"
	"github.com/mgla96/ssh-watcher/internal/notifier"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Starting watcher")
	config, err := config.New()
	if err != nil {
		panic(err)
	}

	notifier := notifier.NewSlackNotifier(config.Slack.WebhookUrl, config.Slack.Channel, config.Slack.Username, config.Slack.Icon, log.Logger)
	processedLineTracker := linetracker.NewFileProcessedLineTracker(config.StateFilePath)

	fileOps := file.FileOps{}
	watcher := app.New(
		config.WatchSettings.LogFileLocation,
		notifier,
		config.HostMachineName,
		config.WatchSettings,
		processedLineTracker,
		fileOps,
	)

	log.Info().Msg(fmt.Sprintf("starting watcher, webhook url: %s, logfile: %s", config.Slack.WebhookUrl, config.WatchSettings.LogFileLocation))
	if err = watcher.Watch(); err != nil {
		panic(err)
	}
}
