package main

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
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
}

func loadConfig() (*Config, error) {
	c := Config{}
	err := envconfig.Process("", &c)
	if err != nil {
		return nil, err
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
	)

	log.Info().Msg(fmt.Sprintf("starting watcher, webhook url: %s, logfile: %s", config.Slack.WebhookUrl, config.LogFileLocation))
	watcher.Watch()
}
