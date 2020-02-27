package util

import "time"

const DateLayout = "2006-01-02"

func TodayStartTime() time.Time {
	timeStr := time.Now().Format(DateLayout)
	t, _ := time.Parse(DateLayout, timeStr)
	return t
}
