package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	defaultContextTimeout time.Duration = 5 * time.Second
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewSlackNotifier(webhookURL, slackChannel, slackUsername, slackIcon string) SlackNotifier {
	return SlackNotifier{
		WebhookURL:    webhookURL,
		SlackChannel:  slackChannel,
		SlackUsername: slackUsername,
		SlackIcon:     slackIcon,
		HttpClient:    &http.Client{},
	}
}

type SlackNotifier struct {
	WebhookURL    string
	SlackChannel  string
	SlackUsername string
	SlackIcon     string
	HttpClient    HTTPClient
}

func (s SlackNotifier) Notify(logLine LogLine) error {
	log.Info().Msg(fmt.Sprintf("Sending notification to slack: User %s %s from IP %s at %s\n", logLine.Username, logLine.EventType, logLine.IpAddress, logLine.LoginTime))

	payloadJson, err := json.Marshal(logLine)
	if err != nil {
		return fmt.Errorf("failed to marshal log line: %w", err)
	}

	slackPayload := SlackPayload{
		Channel:   s.SlackChannel,
		Username:  s.SlackUsername,
		IconEmoji: s.SlackIcon,
		Text:      string(payloadJson),
	}

	log.Info().Msg(fmt.Sprintf("payload: %v", slackPayload))

	payloadJSON, err := json.Marshal(slackPayload)
	if err != nil {
		return fmt.Errorf("error marshaling Slack payload %v: %w", slackPayload, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultContextTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", s.WebhookURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("error creating Slack request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending Slack request: %w", err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading Slack response body: %w", err)
	}
	log.Info().Msg(string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		log.Error().Msg(fmt.Sprintf("response status code not ok: %v", resp.StatusCode))
		return fmt.Errorf("error sending Slack message: %s", resp.Status)
	}

	return nil
}
