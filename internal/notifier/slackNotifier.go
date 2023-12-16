package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type SlackNotifier struct {
	WebhookURL    string
	SlackChannel  string
	SlackUsername string
	SlackIcon     string
}

func (s SlackNotifier) Notify(username, ipAddress, loginTime, eventType, hostMachine string) error {
	log.Info().Msg(fmt.Sprintf("Sending notification to slack: User %s %s from IP %s at %s\n", username, eventType, ipAddress, loginTime))
	//SlackMessagePayload
	payload := LogLine{
		Username:    username,
		IpAddress:   ipAddress,
		LoginTime:   loginTime,
		EventType:   eventType,
		HostMachine: hostMachine,
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	slackPayload := SlackPayload{
		Channel:   s.SlackChannel,
		Username:  s.SlackUsername,
		IconEmoji: s.SlackIcon,
		Text:      string(payloadJson),
	}

	log.Info().Msg(fmt.Sprintf("payload: %v", slackPayload))
	payloadJSON, err := json.Marshal(slackPayload)
	log.Info().Msg(fmt.Sprintf("%s", payloadJSON))
	if err != nil {
		return fmt.Errorf("Error marshaling Slack payload: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", s.WebhookURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("Error creating Slack request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending Slack request: %v", err)
	}
	log.Info().Msg("c")
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	bodyString := string(bodyBytes)
	log.Info().Msg(bodyString)

	if resp.StatusCode != http.StatusOK {
		log.Info().Msg(fmt.Sprintf("response status code not ok: %v", resp.StatusCode))
		return fmt.Errorf("Error sending Slack message: %s", resp.Status)
	}

	return nil
}
