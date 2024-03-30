package notifier

type EventType string

const (
	LoggedIn                          EventType = "logged in"
	FailedLoginAttempt                EventType = "failed login attempt"
	FailedLoginAttemptInvalidUsername EventType = "failed login attempt with invalid username"
)

type LogLine struct {
	Username    string    `json:"username"`
	IpAddress   string    `json:"ip_address"`
	LoginTime   string    `json:"login_time"`
	EventType   EventType `json:"event_type"`
	HostMachine string    `json:"host_machine"`
}

type SlackPayload struct {
	Channel   string `json:"channel"`
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	Text      string `json:"text"`
}
