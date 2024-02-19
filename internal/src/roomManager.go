package src

import (
	"context"
	"fmt"
	"my_app/internal/utils"
	"sync"
	"sync/atomic"
	"time"
)

type RoomManager struct {
	Rooms       *sync.Map
	RoomCounter int32
	mu          sync.Mutex
	context     context.Context
	cancel      context.CancelFunc
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
		room = value.(*Room)
		if success = room.EnterRoom(ctx.Player); success {
			return false
		}
		return true
	})
	for success == false {
		atomic.AddInt32(&m.RoomCounter, 1)
		room = NewRoom()
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

// ============================接口=====================================
// 获取玩家存在的房间
func (rm *RoomManager) GetRoom(ctx *Ctx) (bool, *Room) {
	r, ok := rm.Rooms.Load(ctx.Player.Table.RoomID)
	if !ok || r.(*Room).InRoom(ctx.User.ID) == false {
		ctx.Player.Reset()
		return false, nil
	}
	return true, r.(*Room)
}

// 进入房间
func ReqEnterRoom(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	var room *Room
	if ctx.Player.Table.RoomID > 0 {
		r, ok := ctx.RoomManager.Rooms.Load(ctx.Player.Table.RoomID)
		if !ok || r.(*Room).InRoom(ctx.User.ID) == false {
			ctx.Player.Reset()
		} else {
			room = r.(*Room)
		}
	}
	if room == nil {
		room = ctx.RoomManager.EnterRoom(ctx)

		// 提示其他玩家 有新玩家进入
		ctx.ZClient.Send(utils.Dict{
			"cmd": "ReqZEnterRoom",
			"data": utils.Dict{
				"form_uid":    ctx.User.ID,
				"to_uid_list": room.PlayerIds(ctx.User.ID),
				"message": utils.Dict{
					"uid":  ctx.User.ID,
					"role": ctx.Player.Table.Role,
					"coin": ctx.User.Coin,
					"name": ctx.User.Name,
				},
			},
		})
	}
	ret["room"] = room.GetRet(ctx.Player)
	return ret
}

// 玩家准备
func ReqRoomReady(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	is_ready := data["is_ready"].(bool)
	fmt.Printf("ctx.Player.Table.ID: %v\n", ctx.Player.Table.ID)
	ok, room := ctx.RoomManager.GetRoom(ctx)
	if !ok {
		ret["error"] = "not in room"
		return ret
	}
	if room.ReadyCheck() {
		ret["error"] = "game starting"
		return ret
	}

	ctx.Player.SetReady(is_ready)

	if room.ReadyCheck() {
		room.PlayInit()
	}
	// 提示其他玩家 我做好准备了
	ctx.ZClient.Send(utils.Dict{
		"cmd": "ReqZRoomReady",
		"data": utils.Dict{
			"form_uid":    ctx.User.ID,
			"to_uid_list": room.PlayerIds(ctx.User.ID),
			"message": utils.Dict{
				"is_ready": is_ready,
			},
		},
	})
	ret["is_ready"] = is_ready
	return ret
}

// 看牌
func ReqWatchCards(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	ok, _ := ctx.RoomManager.GetRoom(ctx)
	if !ok {
		ret["error"] = "not in room"
		return ret
	}
	ret["cards"] = ctx.Player.Cards
	return ret
}

// 叫分
func ReqCallScore(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	score := int(data["score"].(float64))
	if score < 1 || score > 3 {
		ret["error"] = "score must be between 1 and 3"
		return
	}
	ok, r := ctx.RoomManager.GetRoom(ctx)
	if !ok {
		ret["error"] = "not in room"
		return ret
	}
	r.CallScore(ctx.Player, score)
	callEnd := score == 3 || r.CallScoreNum == 3
	if callEnd {
		r.ConfirmRole()
	} else {
		r.CallDeskID = (r.CallDeskID + 1) % 3
	}
	// 提示其他玩家 我的叫分
	ctx.ZClient.Send(utils.Dict{
		"cmd": "ReqZCallScore",
		"data": utils.Dict{
			"form_uid":    ctx.User.ID,
			"to_uid_list": r.PlayerIds(ctx.User.ID),
			"message": utils.Dict{
				"score":    score,
				"call_end": callEnd,
			},
		},
	})
	ret["call_end"] = callEnd
	return ret
}

// 获取身份
func ReqGetRole(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	ok, _ := ctx.RoomManager.GetRoom(ctx)
	if !ok {
		ret["error"] = "not in room"
		return ret
	}
	ret["role"] = ctx.Player.Table.Role
	return ret
}

// 出牌
func ReqPlayCards(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	cards := []Card{}
	for _, v := range data["cards"].([]utils.Dict) {
		cards = append(cards, Card{
			Suit:  v["Suit"].(string),
			Value: v["Value"].(string),
		})
	}
	ok, r := ctx.RoomManager.GetRoom(ctx)
	if !ok {
		ret["message"] = "not in room"
		return ret
	}
	r.PlayCard(ctx.Player, cards)
	ret["cards"] = ctx.Player.Cards

	// 提示其他玩家 我的出牌
	ctx.ZClient.Send(utils.Dict{
		"cmd": "ReqZPlayCards",
		"data": utils.Dict{
			"form_uid":    ctx.User.ID,
			"to_uid_list": r.PlayerIds(ctx.User.ID),
			"message": utils.Dict{
				"cards":    cards,
				"card_num": len(ctx.Player.Cards),
			},
		},
	})

	return ret
}
