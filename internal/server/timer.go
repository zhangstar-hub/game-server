package server

import (
	"fmt"
	"my_app/internal/src"
	"time"
)

// 用户在线检测
func (s *Server) UserActiveListener() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	timeoutLimit := 60 * time.Second

	for range ticker.C {
		if s.CloseFlag {
			return
		}
		var canDeleteUser []*src.Ctx
		s.CtxMap.Range(func(key, value interface{}) bool {
			v := value.(*src.Ctx)
			fmt.Printf("time.Since(v.LastActiveTime): %v\n", time.Since(v.LastActiveTime))
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
func (s *Server) AutoSave() {
	saveTime := 10 * time.Minute

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if s.CloseFlag {
			return
		}
		s.CtxMap.Range(func(key, value interface{}) bool {
			v := value.(*src.Ctx)
			if time.Since(v.LastSaveTime) > saveTime {
				v.SaveAll()
			}
			return true
		})
	}
}
