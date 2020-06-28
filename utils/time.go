package utils

import (
	"strings"
	"time"
)

const (
	TimeTpl            = "2006-01-02"
	TimeTplNano        = "2006-01-02 15:04:05.999999999"
	TimeTplRFC3339Nano = "2006-01-02T15:04:05.999999999Z"
)

var tz *time.Location

func LocalTz() *time.Location {
	if tz == nil {
		tz, _ = time.LoadLocation("Asia/Shanghai")
	}
	return tz
}

func CustomDate(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}

	// 如果已经是东八区了,则不转
	//if strings.Contains(t.String(), "+0800 CST") {
	//	return t
	//}

	return t.Add(-time.Hour * 8)
}

// 转换日期成 2006-01-02 15:04:05.999999999的字符串格式
func TimeToString(t time.Time) string {
	s := strings.Replace(t.String(), " +0000 UTC", "", -1)
	return strings.Replace(s, " +0800 CST", "", -1)
}

func BeforeDay() (time.Time, time.Time) {
	t := time.Now()
	now := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, LocalTz())
	before := time.Date(t.Year(), t.Month(), t.Add(-time.Hour * 24).Day(), 0, 0, 0, 0, LocalTz())
	return before, now
}
