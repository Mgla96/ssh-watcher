package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// ServicePrefix - app specific env vars have this prefix.
const ServicePrefix = "WR"

type Slack struct {
	WebhookUrl string `split_words:"true" required:"true"`
	Channel    string `split_words:"true"  default:"#ssh-alerts"`
	Username   string `split_words:"true"  default:"poe-ssh-bot"`
	Icon       string `split_words:"true"  default:":ghost:"`
}

type Config struct {
	HostMachineName string `split_words:"true" required:"true"`
	Slack           *Slack
	WatchSettings   WatchSettings `split_words:"true"`
	// StateFilePath is location of file that keeps track of the last processed line
	// by ssh watcher so restarts of the service do not reprocess all ssh history.
	StateFilePath string `split_words:"true" default:"/var/lib/ssh-watcher/authlog-state"`
}

type WatchSettings struct {
	// AcceptedLogins is a flag to watch for successful logins
	AcceptedLogins bool `default:"true" split_words:"true"`
	// FailedLogins is a flag to watch for failed logins
	FailedLogins bool `default:"true" split_words:"true"`
	// FailedLoginInvalidUsername is a flag to watch for failed logins with invalid username
	FailedLoginInvalidUsername bool `default:"true" split_words:"true"`
	// SleepInterval is the interval in seconds to sleep between log file reads
	SleepInterval int `default:"2" split_words:"true"`
	// LogFileLocation is the location of the log file to watch
	LogFileLocation string `default:"/var/log/auth.log" split_words:"true"`
}

func New() (*Config, error) {
	cfg := Config{}
	err := envconfig.Process(ServicePrefix, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed processing config: %w", err)
	}
	return &cfg, nil
}
