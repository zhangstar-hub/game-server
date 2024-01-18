package utils

import (
	"fmt"
	"time"
)

// 获取今天刷新时间
func TodayFlushTime() time.Time {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return time.Now()
	}
	today := time.Now().In(loc)
	midnight := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, loc)
	return midnight
}

// 获取明天刷新时间
func TomorrowFlushTime() time.Time {
	return TodayFlushTime().Add(24 * time.Hour)
}

// 获取本周一刷新时间
func MondayFlushTime() time.Time {
	today := TodayFlushTime()
	weekday := today.Weekday()
	daysUntilMonday := int(weekday - time.Monday)
	if daysUntilMonday < 0 {
		daysUntilMonday += 7
	}
	mondayMidnight := today.Add(-time.Duration(daysUntilMonday) * 24 * time.Hour)
	return mondayMidnight
}

// 获取下周一刷新时间
func NextMondayFlushTime() time.Time {
	return MondayFlushTime().Add(7 * 24 * time.Hour)
}
