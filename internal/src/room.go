package src

import (
	"errors"
	"fmt"
	"math/rand"
	"my_app/internal/utils"
	"sort"
	"sync"
)

// 房间
type Room struct {
	ID             uint32       // 房间ID
	ZClient        ZMQInterface // 广播器
	BeforeCards    []Card       // 上一次出牌
	BeforePlayDesk int          // 上一次出牌的位置
	Players        [3]*Player   // 玩家
	Score          int64        // 分数
	Mutil          int          // 倍数
	mu             sync.RWMutex // 一个锁
	GameStatus     int          // 游戏状态 0：准备 1:叫分 2:进行
	IsFull         bool         // 是否满房了
	IsClosed       bool         // 房间是否关闭
	IsSpring       bool         // 是否是春天
	winRole        int          // 胜利的角色
	CallDeskID     int          // 出手座位号
	CallScoreNum   int          // 叫分次数
	MaxCallSocre   int          // 叫分最大数
	LastCards      []*Card      // 底牌

	SettleInfo utils.Dict // 结算信息
}

func NewRoom(ZClient ZMQInterface) *Room {
	return &Room{
		Players:    [3]*Player{},
		Score:      1,
		mu:         sync.RWMutex{},
		GameStatus: 0,
		IsFull:     false,
		IsClosed:   false,
		IsSpring:   true,
		winRole:    0,
		ZClient:    ZClient,
		SettleInfo: utils.Dict{},
	}
}

// 房间人数
func (r *Room) PlayerNum() int {
	num := 0
	for _, v := range r.Players {
		if v != nil {
			num++
		}
	}
	return num
}

// 出手流转
func (r *Room) CallConvert() {
	r.CallDeskID = (r.CallDeskID + 1) % 3
}

// 获取房间中玩家ID
func (r *Room) PlayerIds(exculdeID uint) []uint {
	ids := make([]uint, 0)
	for _, v := range r.Players {
		if v == nil {
			continue
		}
		if v.ID == exculdeID {
			continue
		}
		ids = append(ids, v.ID)
	}
	return ids
}

// 玩家进入房间
func (r *Room) EnterRoom(p *Player) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.IsFull {
		return false
	}
	if r.IsClosed {
		return false
	}
	for i := 0; i < len(r.Players); i++ {
		if r.Players[i] == nil {
			p.DeskID = i
			r.Players[i] = p
			break
		}
	}
	if r.PlayerNum() == 3 {
		r.IsFull = true
	}
	return true
}

// 离开房间 主动离开
func (r *Room) LeaveRoom(ctx *Ctx) {
	if r.GameStatus != 0 {
		return
	}

	r.mu.Lock()
	for i, v := range r.Players {
		if v != nil && v.ID == ctx.Player.ID {
			ctx.Player.Reset()
			ctx.User.RoomID = 0
			r.Players[i] = nil
			r.IsFull = false
			break
		}
	}
	r.mu.Unlock()
	r.ZClient.BroastMessage("ReqZLeaveRoom", ctx.Player.ID, r.PlayerIds(ctx.Player.ID), utils.Dict{})
}

// 清理房间
func (r *Room) ClearRoom() {
	if r.GameStatus != 0 {
		return
	}
	r.mu.Lock()
	leaveIds := []uint{}
	aliveIds := []uint{}
	for i, v := range r.Players {
		if IsOnline(v.ID) {
			aliveIds = append(aliveIds, v.ID)
			continue
		}
		leaveIds = append(leaveIds, v.ID)
		r.Players[i] = nil
		r.IsFull = false
	}
	r.mu.Unlock()
	for _, v := range leaveIds {
		r.ZClient.BroastMessage("ReqZLeaveRoom", v, aliveIds, utils.Dict{})
	}

}

// 准备房间中玩家是否全部准备
func (r *Room) ReadyCheck() bool {
	for _, v := range r.Players {
		if v == nil || !v.Ready {
			return false
		}
	}
	return true
}

// 初始化
func (r *Room) StartPlay() {
	r.GameStatus = 1
	r.CallDeskID = rand.Intn(3)
	r.SettleInfo = utils.Dict{}
	// 发牌
	cards := NewCards()
	ShuffleCards(cards)
	for i := 0; i < len(r.Players); i++ {
		r.Players[i].Cards = append([]Card{}, cards[i*17:(i+1)*17]...)
		sort.Slice(r.Players[i].Cards, func(j, k int) bool {
			return r.Players[i].Cards[j].Value > r.Players[i].Cards[k].Value
		})
	}
}

// 洗牌
func ShuffleCards(cards []Card) {
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
}

// 叫分
func (r *Room) CallScore(p *Player, score int) {
	if p.DeskID != r.CallDeskID {
		panic(errors.New("you can't call score"))
	}
	p.Call(score)
	if r.MaxCallSocre < score {
		r.MaxCallSocre = score
	}
	r.CallScoreNum += 1
	if r.CallScoreNum >= 3 || score == 3 {
		r.GameStatus = 2
	}
	r.Mutil = score
}

// 身份确认
func (r *Room) ConfirmRole() {
	for _, v := range r.Players {
		if v.CallScore == r.MaxCallSocre {
			v.ConfirmRole(2)
		} else {
			v.ConfirmRole(1)
		}
	}
}

// 出牌
func (r *Room) PlayCards(p *Player, cards []Card) CardsType {
	if r.CallDeskID != p.DeskID {
		panic(errors.New("not your turn"))
	}
	if r.BeforePlayDesk == p.DeskID && len(cards) == 0 {
		panic(errors.New("your cards are empty"))
	}
	cardsType := GetCardsType(cards)
	if len(cards) > 0 {
		if cardsType == Unknown {
			panic(errors.New("unknown cards type"))
		}
		if r.BeforePlayDesk != p.DeskID && len(r.BeforeCards) > 0 && !IsValidPlay(r.BeforeCards, cards) {
			panic(errors.New("can't play cards"))
		}
		if cardsType == Bomb || cardsType == KingBomb {
			r.Mutil += 2
		}

		p.PlayCards(cards)
		r.BeforePlayDesk = r.CallDeskID
		r.BeforeCards = cards
		if !(r.IsSpring && p.Role == 2) {
			r.IsSpring = false
		}
	}

	r.CallConvert()
	if len(p.Cards) <= 0 {
		r.winRole = p.Role
		if r.IsSpring {
			r.Score *= 2
			r.Mutil += 2
		}
		r.GameOver()
	}
	return cardsType
}

// 结算
func (r *Room) Settle() {
	var CoinPool int64
	playerInfo := utils.Dict{}
	for _, p := range r.Players {
		if r.winRole != p.Role {
			var c int64
			if r.winRole == 1 {
				c = utils.Minint64(GetCoin(p.ID), r.Score)
			} else {
				c = utils.Minint64(GetCoin(p.ID), r.Score*2)
			}
			CoinPool += c
			AddCoin(p.ID, -c)
			playerInfo[fmt.Sprintf("%d", p.ID)] = utils.Dict{
				"winCoins": -c,
			}
		}
	}

	for _, p := range r.Players {
		if r.winRole == p.Role {
			var c int64
			if p.Role == 1 {
				c += CoinPool / 2
			} else {
				c += CoinPool
			}
			AddCoin(p.ID, c)
			playerInfo[fmt.Sprintf("%d", p.ID)] = utils.Dict{
				"winCoins": c,
			}
		}
	}
	r.SettleInfo["player_info"] = playerInfo
	r.SettleInfo["win_role"] = r.winRole
	r.SettleInfo["multi"] = r.Mutil
}

// 游戏结束
func (r *Room) GameOver() {
	r.Settle()
	r.GameStatus = 0
	for _, p := range r.Players {
		p.Reset()
	}
	r.ClearRoom()
}

// 关闭房间
func (r *Room) Close() {
	r.IsClosed = true
}

// 判断一个玩家是否在房间中
func (r *Room) InRoom(uid uint) bool {
	for _, v := range r.Players {
		if v != nil && v.ID == uid {
			return true
		}
	}
	return false
}

// 玩家信息
func (r *Room) PlayersInfo() []utils.Dict {
	players := []utils.Dict{}
	for _, v := range r.Players {
		if v == nil {
			players = append(players, utils.Dict{})
		} else {
			players = append(players, v.GetRet())
		}
	}
	return players
}

func (r *Room) GetRet(p *Player) utils.Dict {
	ret := utils.Dict{
		"uid":          p.ID,
		"room_id":      r.ID,
		"game_status":  r.GameStatus,
		"score":        r.Score,
		"call_desk":    r.CallDeskID,
		"played_cards": CardsToValue(r.BeforeCards),
		"cards":        CardsToValue(p.Cards),
	}
	ret["players"] = r.PlayersInfo()
	return ret
}
