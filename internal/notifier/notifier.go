package notifier

type Notifier interface {
	Notify(username, ipAddress, loginTime, eventType, hostMachine string) error
}

type LogLine struct {
	Username    string `json:"username"`
	IpAddress   string `json:"ip_address"`
	LoginTime   string `json:"login_time"`
	EventType   string `json:"event_type"`
	HostMachine string `json:"host_machine"`
}

type SlackPayload struct {
	Channel   string `json:"channel"`
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	Text      string `json:"text"`
}
