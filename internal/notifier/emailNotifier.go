package notifier

type EmailNotifier struct{}

func (e EmailNotifier) Notify(logLine LogLine) error {
	panic("not implemented")
}
