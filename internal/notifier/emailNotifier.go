package notifier

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type EmailNotifier struct{}

func (e EmailNotifier) Notify(username, ipAddress, loginTime, eventType, hostMachine string) error {
	log.Info().Msg(fmt.Sprintf("Sending notification to email: User %s %s from IP %s at %s\n", username, eventType, ipAddress, loginTime))
	panic("not implemented")
}
