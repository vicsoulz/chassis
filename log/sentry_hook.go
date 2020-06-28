package log

import (
	"errors"
	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

type SentryHook struct {
	DSN   string
	Level []logrus.Level
}

func (s *SentryHook) Levels() []logrus.Level {
	if s.Level == nil || len(s.Level) == 0 {
		return []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		}
	}
	return s.Level
}

// trace init stack
func (s SentryHook) trace() *raven.Stacktrace {
	return raven.NewStacktrace(0, 2, nil)
}

func (s *SentryHook) Fire(e *logrus.Entry) error {
	client, err := raven.New(s.DSN)
	if err != nil {
		return err
	}

	packet := &raven.Packet{
		Message:    e.Message,
		Interfaces: []raven.Interface{raven.NewException(errors.New(e.Level.String()), s.trace())}}
	_, ch := client.Capture(packet, nil)

	if err = <-ch; err != nil {
		return err
	}

	return nil
}
