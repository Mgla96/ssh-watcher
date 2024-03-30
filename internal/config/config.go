package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// ServicePrefix - app specific env vars have this prefix.
const ServicePrefix = "WR"

type SlackConfig struct {
	WebhookUrl string `envconfig:"SLACK_WEBHOOK_URL"`
	Channel    string `envconfig:"SLACK_CHANNEL" default:"#ssh-alerts"`
	Username   string `envconfig:"SLACK_USERNAME" defaul:"poe-ssh-bot"`
	Icon       string `envconfig:"SLACK_ICON" default:":ghost:"`
}

type Config struct {
	HostMachineName string       `envconfig:"HOST_MACHINE_NAME" required:"true"`
	Slack           *SlackConfig `envconfig:"SLACK"`
	LogFileLocation string       `envconfig:"WATCH_LOGFILE" default:"/var/log/auth.log"`
	// WatchAcceptedLogin              bool         `envconfig:"WATCH_SETTINGS_ACCEPTED_LOGIN" default:"true"`
	// WatchFailedLogin                bool         `envconfig:"WATCH_SETTINGS_FAILED_LOGIN" default:"true"`
	// WatchFailedLoginInvalidUsername bool         `envconfig:"WATCH_SETTINGS_FAILED_LOGIN_INVALID_USERNAME" default:"false"`
	// WatchSleepIntervalSeconds       int          `envconfig:"WATCH_SETTINGS_SLEEP_INTERVAL_SECONDS" default:"2"`
	// StateFilePath is location of file that keeps track of the last processed line
	// by ssh watcher so restarts of the service do not reprocess all ssh history.
	WatchSettings WatchSettings `split_words:"true"`
	StateFilePath string        `envconfig:"STATE_FILE_PATH" default:"/var/lib/ssh-watcher/authlog-state"`
}

type WatchSettings struct {
	AcceptedLogins             bool `default:"true" split_words:"true"`
	FailedLogins               bool `default:"false" split_words:"true"`
	FailedLoginInvalidUsername bool `default:"true" split_words:"true"`
	SleepInterval              int  `default:"2" split_words:"true"`
}

func New() (*Config, error) {
	cfg := Config{}
	err := envconfig.Process(ServicePrefix, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed processing config: %w", err)
	}
	return &cfg, nil
}
