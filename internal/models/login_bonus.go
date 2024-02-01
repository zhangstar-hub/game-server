package models

import (
	"time"
)

type LoginBonusModel struct {
	ID              uint      `gorm:"primary_key"`
	SiginTimes      int       // 签到次数
	DailyFlushTime  time.Time // 日刷新时间
	WeeklyFlushTime time.Time // 周刷新时间
}
