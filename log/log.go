package log

import (
	"errors"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Log       *logrus.Logger
	SentryDSN string
}

var (
	l    *logrus.Logger
	conf *Config
)

func Init(c *Config) error {
	if c == nil {
		return errors.New("config is nil")
	}
	conf = c

	if c.Log == nil {
		l = logrus.New()
	} else {
		l = c.Log
	}

	initHook()
	return nil
}

func initHook() {
	ch := &CallerHook{}
	if l.Level == logrus.DebugLevel {
		ch.Level = logrus.AllLevels
	}
	l.AddHook(ch)

	if conf.SentryDSN != "" {
		l.AddHook(&SentryHook{
			DSN: conf.SentryDSN,
		})
	}
}

func Get() *logrus.Logger {
	return l
}

func Trace(args ...interface{}) {
	l.Trace(args)
}

func Tracef(format string, args ...interface{}) {
	l.Tracef(format, args...)
}

func Debug(args ...interface{}) {
	l.Debug(args)
}

func Debugf(format string, args ...interface{}) {
	l.Debugf(format, args...)
}

func Info(args ...interface{}) {
	l.Info(args)
}

func Infof(format string, args ...interface{}) {
	l.Infof(format, args...)
}

func Warn(args ...interface{}) {
	l.Warn(args)
}

func Warnf(format string, args ...interface{}) {
	l.Warnf(format, args...)
}

func Error(args ...interface{}) {
	l.Error(args)
}

func Errorf(format string, args ...interface{}) {
	l.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	l.Fatal(args)
}

func Fatalf(format string, args ...interface{}) {
	l.Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	l.Panic(args)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return l.WithFields(fields)
}
