package log

import (
	"errors"
	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Log       *logrus.Logger
	SentryDSN string
}

var (
	Log  *logrus.Logger
	conf *Config
)

func Init(c *Config) (*logrus.Logger, error) {
	if c == nil {
		return nil, errors.New("config is nil")
	}
	conf = c

	if c.Log == nil {
		Log = logrus.New()
	} else {
		Log = c.Log
	}

	err := initHook()
	return Log, err
}

func initHook() error {
	Log.SetReportCaller(true)
	if conf.SentryDSN != "" {
		hook, err := logrus_sentry.NewSentryHook(conf.SentryDSN, []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		})

		if err != nil {
			return err
		}
		Log.Hooks.Add(hook)
		hook.StacktraceConfiguration.Enable = true
	}
	return nil
}
