package server

import (
	"my_app/internal/src"
	"time"
)

// 用户在线检测
func UserActiveListener() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	timeoutLimit := 60 * time.Second

	for range ticker.C {
		var canDeleteUser []*src.Ctx
		for _, v := range src.Users {
			if time.Since(v.LastActiveTime) > timeoutLimit {
				canDeleteUser = append(canDeleteUser, v)
			}
		}
		for _, v := range canDeleteUser {
			v.Close()
		}
	}
}
