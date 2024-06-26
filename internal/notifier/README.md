<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# notifier

```go
import "github.com/mgla96/ssh-watcher/internal/notifier"
```

## Index

- [type EmailNotifier](<#EmailNotifier>)
  - [func \(e EmailNotifier\) Notify\(logLine LogLine\) error](<#EmailNotifier.Notify>)
- [type EventType](<#EventType>)
- [type HTTPClient](<#HTTPClient>)
- [type LogLine](<#LogLine>)
- [type SlackNotifier](<#SlackNotifier>)
  - [func NewSlackNotifier\(webhookURL, slackChannel, slackUsername, slackIcon string\) SlackNotifier](<#NewSlackNotifier>)
  - [func \(s SlackNotifier\) Notify\(logLine LogLine\) error](<#SlackNotifier.Notify>)
- [type SlackPayload](<#SlackPayload>)


<a name="EmailNotifier"></a>
## type [EmailNotifier](<https://github.com/Mgla96/ssh-watcher/blob/main/internal/notifier/email_notifier.go#L3>)



```go
type EmailNotifier struct{}
```

<a name="EmailNotifier.Notify"></a>
### func \(EmailNotifier\) [Notify](<https://github.com/Mgla96/ssh-watcher/blob/main/internal/notifier/email_notifier.go#L5>)

```go
func (e EmailNotifier) Notify(logLine LogLine) error
```



<a name="EventType"></a>
## type [EventType](<https://github.com/Mgla96/ssh-watcher/blob/main/internal/notifier/notifier.go#L3>)



```go
type EventType string
```

<a name="LoggedIn"></a>

```go
const (
    LoggedIn                          EventType = "logged in"
    FailedLoginAttempt                EventType = "failed login attempt"
    FailedLoginAttemptInvalidUsername EventType = "failed login attempt with invalid username"
)
```

<a name="HTTPClient"></a>
## type [HTTPClient](<https://github.com/Mgla96/ssh-watcher/blob/main/internal/notifier/slack_notifier.go#L19-L21>)



```go
type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}
```

<a name="LogLine"></a>
## type [LogLine](<https://github.com/Mgla96/ssh-watcher/blob/main/internal/notifier/notifier.go#L11-L17>)



```go
type LogLine struct {
    Username    string    `json:"username"`
    IpAddress   string    `json:"ip_address"`
    LoginTime   string    `json:"login_time"`
    EventType   EventType `json:"event_type"`
    HostMachine string    `json:"host_machine"`
}
```

<a name="SlackNotifier"></a>
## type [SlackNotifier](<https://github.com/Mgla96/ssh-watcher/blob/main/internal/notifier/slack_notifier.go#L33-L39>)



```go
type SlackNotifier struct {
    WebhookURL    string
    SlackChannel  string
    SlackUsername string
    SlackIcon     string
    HttpClient    HTTPClient
}
```

<a name="NewSlackNotifier"></a>
### func [NewSlackNotifier](<https://github.com/Mgla96/ssh-watcher/blob/main/internal/notifier/slack_notifier.go#L23>)

```go
func NewSlackNotifier(webhookURL, slackChannel, slackUsername, slackIcon string) SlackNotifier
```



<a name="SlackNotifier.Notify"></a>
### func \(SlackNotifier\) [Notify](<https://github.com/Mgla96/ssh-watcher/blob/main/internal/notifier/slack_notifier.go#L41>)

```go
func (s SlackNotifier) Notify(logLine LogLine) error
```



<a name="SlackPayload"></a>
## type [SlackPayload](<https://github.com/Mgla96/ssh-watcher/blob/main/internal/notifier/notifier.go#L19-L24>)



```go
type SlackPayload struct {
    Channel   string `json:"channel"`
    Username  string `json:"username"`
    IconEmoji string `json:"icon_emoji"`
    Text      string `json:"text"`
}
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
