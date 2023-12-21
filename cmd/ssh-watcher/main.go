package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mgla96/ssh-watcher/internal/notifier"
	"github.com/mgla96/ssh-watcher/internal/watcher"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	defaultSlackChannel                           = "#ssh-alerts"
	defaultSlackUsername                          = "poe-ssh-bot"
	defaultSlackIcon                              = ":ghost:"
	defaultLogFileLocation                        = "/var/log/auth.log"
	defaultWatchAcceptedLogins                    = true
	defaultWatchFailedLogins                      = false
	defaultWatchFailedLoginInvalidUsername        = false
	defaultWatchSleepIntervalSeconds              = 2
	envKeyHostUrl                                 = "HOST_URL"
	envKeySlackWebhookUrl                         = "SLACK_WEBHOOK_URL"
	envKeySlackChannel                            = "SLACK_CHANNEL"
	envKeySlackUsername                           = "SLACK_USERNAME"
	envKeySlackIcon                               = "SLACK_ICON"
	envKeyWatchLogfile                            = "WATCH_LOGFILE"
	envKeyWatchSettingsAcceptedLogin              = "WATCH_SETTINGS_ACCEPTED_LOGIN"
	envKeyWatchSettingsFailedLogin                = "WATCH_SETTINGS_FAILED_LOGIN"
	envKeyWatchSettingsFailedLoginInvalidUsername = "WATCH_SETTINGS_FAILED_LOGIN_INVALID_USERNAME"
	envKeyWatchInterval                           = "WATCH_INTERVAL"
)

type Config struct {
	HostUrl                         string
	WebhookUrl                      string
	SlackChannel                    string
	SlackUsername                   string
	SlackIcon                       string
	LogFileLocation                 string
	WatchAcceptedLogin              bool
	WatchFailedLogin                bool
	WatchFailedLoginInvalidUsername bool
	WatchSleepIntervalSeconds       int
}

func loadConfig() Config {
	c := Config{}

	c.HostUrl = getEnvOrPanic(envKeyHostUrl)
	c.WebhookUrl = getEnvOrPanic(envKeySlackWebhookUrl)

	c.SlackChannel = getEnvOrDefault(envKeySlackChannel, defaultSlackChannel)
	c.SlackUsername = getEnvOrDefault(envKeySlackUsername, defaultSlackUsername)
	c.SlackIcon = getEnvOrDefault(envKeySlackIcon, defaultSlackIcon)
	c.LogFileLocation = getEnvOrDefault(envKeyWatchLogfile, defaultLogFileLocation)

	c.WatchAcceptedLogin = parseBoolEnv(envKeyWatchSettingsAcceptedLogin, defaultWatchAcceptedLogins)
	c.WatchFailedLogin = parseBoolEnv(envKeyWatchSettingsFailedLogin, defaultWatchFailedLogins)
	c.WatchFailedLoginInvalidUsername = parseBoolEnv(envKeyWatchSettingsFailedLoginInvalidUsername, defaultWatchFailedLoginInvalidUsername)

	c.WatchSleepIntervalSeconds = parseIntEnv(envKeyWatchInterval, defaultWatchSleepIntervalSeconds)

	return c
}

func getEnvOrPanic(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Sprintf("%s not set", key))
	}
	return value
}

func getEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Warn().Msgf("%s not set, defaulting to %s", key, defaultValue)
		return defaultValue
	}
	return value
}

func parseBoolEnv(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Warn().Msgf("%s not parsable, defaulting to %t", key, defaultValue)
		return defaultValue
	}
	return value
}

func parseIntEnv(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	value, err := strconv.ParseInt(valueStr, 10, 0)
	if err != nil {
		log.Warn().Msgf("%s not parsable, defaulting to %d", key, defaultValue)
		return defaultValue
	}
	return int(value)
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Starting watcher")
	config := loadConfig()

	notifier := notifier.NewSlackNotifier(config.WebhookUrl, config.SlackChannel, config.SlackUsername, config.SlackIcon)
	watcher := watcher.NewLogWatcher(
		config.LogFileLocation,
		notifier,
		config.HostUrl,
		watcher.WatchSettings{
			WatchAcceptedLogins:             config.WatchAcceptedLogin,
			WatchFailedLogins:               config.WatchFailedLogin,
			WatchFailedLoginInvalidUsername: config.WatchFailedLoginInvalidUsername,
			WatchSleepInterval:              time.Duration(config.WatchSleepIntervalSeconds) * time.Second,
		},
	)

	log.Info().Msg(fmt.Sprintf("starting watcher, webhook url: %s, logfile: %s", config.WebhookUrl, config.LogFileLocation))
	watcher.Watch()
}
