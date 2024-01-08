package main

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/mgla96/ssh-watcher/internal/linetracker"
	"github.com/mgla96/ssh-watcher/internal/notifier"
	"github.com/mgla96/ssh-watcher/internal/watcher"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type SlackConfig struct {
	WebhookUrl string `envconfig:"SLACK_WEBHOOK_URL"`
	Channel    string `envconfig:"SLACK_CHANNEL" default:"#ssh-alerts"`
	Username   string `envconfig:"SLACK_USERNAME" defaul:"poe-ssh-bot"`
	Icon       string `envconfig:"SLACK_ICON" default:":ghost:"`
}

type Config struct {
	HostMachineName                 string       `envconfig:"HOST_MACHINE_NAME" required:"true"`
	Slack                           *SlackConfig `envconfig:"SLACK"`
	LogFileLocation                 string       `envconfig:"WATCH_LOGFILE" default:"/var/log/auth.log"`
	WatchAcceptedLogin              bool         `envconfig:"WATCH_SETTINGS_ACCEPTED_LOGIN" default:"true"`
	WatchFailedLogin                bool         `envconfig:"WATCH_SETTINGS_FAILED_LOGIN" default:"true"`
	WatchFailedLoginInvalidUsername bool         `envconfig:"WATCH_SETTINGS_FAILED_LOGIN_INVALID_USERNAME" default:"false"`
	WatchSleepIntervalSeconds       int          `envconfig:"WATCH_SETTINGS_SLEEP_INTERVAL_SECONDS" default:"2"`
	// StateFilePath is location of file that keeps track of the last processed line
	// by ssh watcher so restarts of the service do not reprocess all ssh history.
	StateFilePath string `envconfig:"STATE_FILE_PATH default:"/var/lib/ssh-watcher/authlog-state"`
}

func loadConfig() (*Config, error) {
	c := Config{}
	err := envconfig.Process("", &c)
	if err != nil {
		return nil, fmt.Errorf("failed processing config: %w", err)
	}
	return &c, nil
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Starting watcher")
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	notifier := notifier.NewSlackNotifier(config.Slack.WebhookUrl, config.Slack.Channel, config.Slack.Username, config.Slack.Icon)
	processedLineTracker := linetracker.NewFileProcessedLineTracker(config.StateFilePath)
	watcher := watcher.NewLogWatcher(
		config.LogFileLocation,
		notifier,
		config.HostMachineName,
		watcher.WatchSettings{
			WatchAcceptedLogins:             config.WatchAcceptedLogin,
			WatchFailedLogins:               config.WatchFailedLogin,
			WatchFailedLoginInvalidUsername: config.WatchFailedLoginInvalidUsername,
			WatchSleepInterval:              time.Duration(config.WatchSleepIntervalSeconds) * time.Second,
		},
		processedLineTracker,
	)

	log.Info().Msg(fmt.Sprintf("starting watcher, webhook url: %s, logfile: %s", config.Slack.WebhookUrl, config.LogFileLocation))
	if err = watcher.Watch(); err != nil {
		panic(err)
	}
}
