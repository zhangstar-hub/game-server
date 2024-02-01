package src

import (
	"context"
	"my_app/internal/utils"
	"sync"
	"sync/atomic"
	"time"
)

type RoomManager struct {
	Rooms        *sync.Map
	RoomCounter  int32
	mu           sync.Mutex
	context      context.Context
	cancel       context.CancelFunc
	minRoomCount int // 最小保留的房间数量
}

func NewRoomManager() *RoomManager {
	r_manager := &RoomManager{
		Rooms:       new(sync.Map),
		RoomCounter: 0,
		mu:          sync.Mutex{},
	}
	r_manager.context, r_manager.cancel = context.WithCancel(context.Background())
	go r_manager.AutoClearRoom()
	return r_manager
}

// 进入一个房间
func (m *RoomManager) EnterRoom(ctx *Ctx) (room *Room) {
	success := false
	ctx.Player.Reset()

	m.Rooms.Range(func(key, value any) bool {
		room := value.(*Room)
		if success = room.EnterRoom(ctx.Player); success {
			return false
		}
		return true
	})

	for success == false {
		atomic.AddInt32(&m.RoomCounter, 1)
		room = NewRoom(ctx)
		room.ID = m.RoomCounter
		m.Rooms.Store(m.RoomCounter, room)
		if success = room.EnterRoom(ctx.Player); success {
			break
		}
	}
	return
}

// 定时清理空房间
func (rm *RoomManager) AutoClearRoom() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if utils.SynMapLength(rm.Rooms) <= rm.minRoomCount {
				continue
			}
			closeRooms := []*Room{}
			rm.Rooms.Range(func(key, value any) bool {
				r := value.(*Room)
				if len(r.Players) == 0 {
					r.Close()
					closeRooms = append(closeRooms, r)
				}
				return true
			})
			for _, r := range closeRooms {
				rm.Rooms.Delete(r.ID)
			}
		case <-rm.context.Done():
			return
		}
	}
}

// 关闭房间管理
func (rm *RoomManager) Close() error {
	rm.cancel()
	rm.Rooms.Range(func(key, value any) bool {
		r := value.(*Room)
		r.Close()
		return true
	})
	return nil
}

func ReqEnterRoom(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	room := ctx.RoomManager.EnterRoom(ctx)
	ret["room"] = room.GetRet()
	return ret
}
