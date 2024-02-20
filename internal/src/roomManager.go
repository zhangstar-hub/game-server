package src

import (
	"context"
	"encoding/json"
	"fmt"
	"my_app/internal/utils"
	"os"
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
	for success == false {
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

// 获取玩家存在的房间
func (rm *RoomManager) GetRoom(ctx *Ctx) (bool, *Room) {
	r, ok := rm.Rooms.Load(ctx.User.RoomID)
	if !ok || r.(*Room).InRoom(ctx.User.ID) == false {
		ctx.Player.Reset()
		return false, nil
	}
	return true, r.(*Room)
}

// ============================接口=====================================
// 进入房间
func ReqEnterRoom(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	var room *Room
	if ctx.User.RoomID > 0 {
		r, ok := ctx.RoomManager.Rooms.Load(ctx.User.RoomID)
		if !ok || r.(*Room).InRoom(ctx.User.ID) == false {
			ctx.Player.Reset()
		} else {
			room = r.(*Room)
		}
	}
	if room == nil {
		room = ctx.RoomManager.EnterRoom(ctx)

		// 提示其他玩家 有新玩家进入
		ctx.ZClient.SendMessage(
			"ReqZEnterRoom",
			ctx.User.ID,
			room.PlayerIds(ctx.User.ID),
			utils.Dict{
				"uid":  ctx.User.ID,
				"role": ctx.Player.Role,
				"coin": ctx.User.Coin,
				"name": ctx.User.Name,
			},
		)
	}
	ret["room"] = room.GetRet(ctx.Player)
	return ret
}

// 玩家准备
func ReqRoomReady(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	is_ready := data["is_ready"].(bool)
	fmt.Printf("ctx.Player.Table.ID: %v\n", ctx.User.ID)
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
		room.StartPlay()
	}
	// 提示其他玩家 我做好准备了
	ctx.ZClient.SendMessage(
		"ReqZRoomReady",
		ctx.User.ID,
		room.PlayerIds(ctx.User.ID),
		utils.Dict{
			"is_ready": is_ready,
		},
	)

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
		r.CallConvert()
	}
	// 提示其他玩家 我的叫分
	ctx.ZClient.SendMessage(
		"ReqZCallScore",
		ctx.User.ID,
		r.PlayerIds(ctx.User.ID),
		utils.Dict{
			"score":    score,
			"call_end": callEnd,
		},
	)
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
	ret["role"] = ctx.Player.Role
	return ret
}

// 出牌
func ReqPlayCards(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	cards := []Card{}
	// for _, v := range data["cards"].([]utils.Dict) {
	// 	cards = append(cards, Card{
	// 		Suit:  v["Suit"].(string),
	// 		Value: v["Value"].(string),
	// 	})
	// }

	// temp
	file, err := os.Open("configs/cards.json")
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cards); err != nil {
		fmt.Println("解码 JSON 失败:", err)
		return
	}
	// temp

	ok, r := ctx.RoomManager.GetRoom(ctx)
	if !ok {
		ret["message"] = "not in room"
		return ret
	}
	r.PlayCards(ctx.Player, cards)
	ret["cards"] = ctx.Player.Cards

	// 提示其他玩家 我的出牌
	ctx.ZClient.SendMessage(
		"ReqZPlayCards",
		ctx.User.ID,
		r.PlayerIds(ctx.User.ID),
		utils.Dict{
			"cards":    cards,
			"card_num": len(ctx.Player.Cards),
		},
	)
	return ret
}
