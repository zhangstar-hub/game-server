package src

import (
	"context"
	"my_app/internal/utils"
	"sync"
	"sync/atomic"
	"time"
)

type RoomManager struct {
	Rooms       *sync.Map
	RoomCounter uint32
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
	for !success {
		atomic.AddUint32(&m.RoomCounter, 1)
		room = NewRoom(ctx.ZClient)
		room.ID = uint32(m.RoomCounter)
		m.Rooms.Store(m.RoomCounter, room)
		if success = room.EnterRoom(ctx.Player); success {
			break
		}
	}
	ctx.User.RoomID = uint32(room.ID)
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
				r.mu.Lock()
				if r.PlayerNum() == 0 {
					r.Close()
					closeRooms = append(closeRooms, r)
				}
				r.mu.Unlock()
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

// 获取玩家存在的房间
func (rm *RoomManager) GetRoom(roomID uint32, uid uint) (bool, *Room) {
	r, ok := rm.Rooms.Load(roomID)
	if !ok || !r.(*Room).InRoom(uid) {
		return false, nil
	}
	return true, r.(*Room)
}

// ============================接口=====================================
// 进入房间
func ReqEnterRoom(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	ok, room := ctx.RoomManager.GetRoom(ctx.User.RoomID, ctx.User.ID)
	if !ok {
		ret["error"] = "not in room"
	}
	if room == nil {
		room = ctx.RoomManager.EnterRoom(ctx)
	}

	// 提示其他玩家 有新玩家进入
	ctx.ZClient.BroastMessage(
		"ReqZEnterRoom",
		ctx.User.ID,
		room.PlayerIds(ctx.User.ID),
		ctx.Player.GetRet(),
	)

	ret["room"] = room.GetRet(ctx.Player)
	return ret
}

// 离开房间
func ReqLeaveRoom(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	ok, room := ctx.RoomManager.GetRoom(ctx.User.RoomID, ctx.User.ID)
	if !ok {
		ret["error"] = "not in room"
	}
	room.LeaveRoom(ctx)
	return ret
}

// 玩家准备
func ReqRoomReady(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	is_ready := data["is_ready"].(bool)
	ok, room := ctx.RoomManager.GetRoom(ctx.User.RoomID, ctx.User.ID)
	if !ok {
		ret["error"] = "not in room"
		return ret
	}
	if room.ReadyCheck() {
		ret["error"] = "game starting"
		return ret
	}

	ctx.Player.SetReady(is_ready)

	is_started := room.ReadyCheck()
	if is_started {
		room.StartPlay()
	}
	// 提示其他玩家 我做好准备了
	ctx.ZClient.BroastMessage(
		"ReqZRoomReady",
		ctx.User.ID,
		room.PlayerIds(ctx.User.ID),
		utils.Dict{
			"is_ready":    is_ready,
			"game_status": room.GameStatus,
			"call_desk":   room.CallDeskID,
		},
	)

	ret["is_ready"] = is_ready
	ret["game_status"] = room.GameStatus
	ret["call_desk"] = room.CallDeskID
	return ret
}

// 看牌
func ReqWatchCards(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	ok, room := ctx.RoomManager.GetRoom(ctx.User.RoomID, ctx.User.ID)
	if !ok {
		ret["error"] = "not in room"
		return ret
	}
	players_cards_num := []utils.Dict{}
	for _, v := range room.Players {
		players_cards_num = append(players_cards_num, utils.Dict{
			"id":       v.ID,
			"card_num": len(v.Cards),
		})
	}
	ret["game_status"] = room.GameStatus
	ret["cards"] = CardsToValue(ctx.Player.Cards)
	ret["players_cards_num"] = players_cards_num
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
	ok, r := ctx.RoomManager.GetRoom(ctx.User.RoomID, ctx.User.ID)
	if !ok {
		ret["error"] = "not in room"
		return ret
	}
	if r.GameStatus != 1 {
		ret["error"] = "can't call score"
		return
	}
	if score <= r.MaxCallSocre {
		ret["error"] = "must geater than before call score"
	}
	r.CallScore(ctx.Player, score)
	if score == 3 || r.CallScoreNum == 3 {
		r.ConfirmRole()
	} else {
		r.CallConvert()
	}
	// 提示其他玩家 我的叫分
	ctx.ZClient.BroastMessage(
		"ReqZCallScore",
		ctx.User.ID,
		r.PlayerIds(ctx.User.ID),
		utils.Dict{
			"score":       score,
			"game_status": r.GameStatus,
			"call_desk":   r.CallDeskID,
		},
	)
	ret["call_desk"] = r.CallDeskID
	ret["game_status"] = r.GameStatus
	return ret
}

// 获取身份
func ReqGetRole(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	ok, _ := ctx.RoomManager.GetRoom(ctx.User.RoomID, ctx.User.ID)
	if !ok {
		ret["error"] = "not in room"
		return ret
	}
	ret["role"] = ctx.Player.Role
	return ret
}

// 出牌
func ReqPlayCards(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	cardsValue := []int{}
	for _, v := range data["cards"].([]interface{}) {
		cardsValue = append(cardsValue, int(v.(float64)))
	}
	cards := ValueToCards(cardsValue)
	ok, r := ctx.RoomManager.GetRoom(ctx.User.RoomID, ctx.User.ID)
	if !ok {
		ret["message"] = "not in room"
		return ret
	}

	cardsType := r.PlayCards(ctx.Player, cards)
	game_status := r.GameStatus
	if len(ctx.Player.Cards) == 0 {
		game_status = 3
	}

	ret["cards"] = CardsToValue(ctx.Player.Cards)
	ret["cards_type"] = cardsType
	ret["game_status"] = game_status
	ret["call_desk"] = r.CallDeskID
	ret["played_cards"] = data["cards"]

	// 提示其他玩家 我的出牌
	ctx.ZClient.BroastMessage(
		"ReqZPlayCards",
		ctx.User.ID,
		r.PlayerIds(ctx.User.ID),
		utils.Dict{
			"call_desk":    r.CallDeskID,
			"game_status":  game_status,
			"played_cards": data["cards"],
			"cards_type":   cardsType,
			"card_num":     len(ctx.Player.Cards),
		},
	)
	return ret
}
