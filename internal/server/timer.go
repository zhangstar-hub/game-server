package server

import (
	"my_app/internal/context"
	"time"
)

// 用户在线检测
func UserActiveListener() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	timeoutLimit := 60 * time.Second

	for range ticker.C {
		var canDeleteUser []*context.Ctx
		context.Users.Range(func(key, value interface{}) bool {
			v := value.(*context.Ctx)
			if time.Since(v.LastActiveTime) > timeoutLimit {
				canDeleteUser = append(canDeleteUser, v)
			}
			return true
		})
		for _, v := range canDeleteUser {
			v.Close()
		}
	}
}

// 定期存储数据 防止数据丢失
func AutoSave() {
	saveTime := 10 * time.Minute

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		context.Users.Range(func(key, value interface{}) bool {
			v := value.(*context.Ctx)
			if time.Since(v.LastSaveTime) > saveTime {
				v.SaveAll()
			}
			return true
		})
	}
}
