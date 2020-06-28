package log

import (
	"github.com/sirupsen/logrus"
	"runtime"
	"strconv"
	"strings"
)

type CallerHook struct {
	Level []logrus.Level
}

func (h *CallerHook) Levels() []logrus.Level {
	if h.Level == nil || len(h.Level) == 0 {
		return []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.DebugLevel,
		}
	}
	return h.Level
}

func (h *CallerHook) Fire(e *logrus.Entry) error {
	_, file, line, ok := runtime.Caller(8)

	if ok {
		idx := strings.LastIndex(file, "src")
		if idx >= 0 {
			file = file[idx+4:]
		}

		indexFunc := func(file string) string {
			backup := "/" + file
			lastSlashIndex := strings.LastIndex(backup, "/")
			if lastSlashIndex < 0 {
				return backup
			}
			secondLastSlashIndex := strings.LastIndex(backup[:lastSlashIndex], "/")
			if secondLastSlashIndex < 0 {
				return backup[lastSlashIndex+1:]
			}
			return backup[secondLastSlashIndex+1:]
		}
		file = indexFunc(file) + ":" + strconv.Itoa(line)
	}
	e.Data = map[string]interface{}{
		"file": file,
	}
	return nil
}
