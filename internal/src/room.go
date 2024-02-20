package src

import (
	"errors"
	"math/rand"
	"my_app/internal/utils"
	"sort"
	"sync"
)

// 房间
type Room struct {
	ID             int32        // 房间ID
	BeforeCards    []Card       // 上一次出牌
	BeforePlayDesk int          // 上一次出牌的位置
	Cards          []Card       // 桌面上的纸牌
	Players        [3]*Player   // 玩家
	Score          uint64       // 分数
	CoinPool       uint64       //	奖池
	mu             sync.RWMutex // 一个锁
	IsOver         bool         // 游戏是否结束
	IsFull         bool         // 是否满房了
	IsClosed       bool         // 房间是否关闭
	winRole        int          // 胜利的角色
	CallDeskID     int          // 出手座位号
	CallScoreNum   int          // 叫分次数
	MaxCallSocre   int          // 叫分最大数
}

func NewRoom() *Room {
	return &Room{
		Cards:    NewCards(),
		Players:  [3]*Player{},
		Score:    0,
		CoinPool: 0,
		mu:       sync.RWMutex{},
		IsOver:   false,
		IsFull:   false,
		IsClosed: false,
		winRole:  0,
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
func (r *Room) PlayerIds(exculdeId uint) []uint {
	ids := make([]uint, 0)
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, v := range r.Players {
		if v == nil {
			continue
		}
		if v.Table.ID == exculdeId {
			continue
		}
		ids = append(ids, v.Table.ID)
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
			r.Players[i] = p
			p.Table.RoomID = r.ID
			p.Table.DeskID = i
			break
		}
	}
	if r.PlayerNum() == 3 {
		r.IsFull = true
	}
	return true
}

// 离开房间
func (r *Room) LeaveRoom(p *Player) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, v := range r.Players {
		if v == p {
			r.Players[i] = nil
			r.IsFull = false
			p.Reset()
			break
		}
	}
}

// 准备房间中玩家是否全部准备
func (r *Room) ReadyCheck() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, v := range r.Players {
		if v == nil || !v.Table.Ready {
			return false
		}
	}
	return true
}

// 初始化
func (r *Room) StartPlay() {
	r.CallDeskID = rand.Intn(3)

	// 发牌
	for i := 0; i < len(r.Players); i++ {
		r.Players[i].Cards = append([]Card{}, r.Cards[i*17:(i+1)*17]...)
		sort.Slice(r.Players[i].Cards, func(j, k int) bool {
			return r.Players[i].Cards[j].Value < r.Players[i].Cards[k].Value
		})
	}
}

// 叫分
func (r *Room) CallScore(p *Player, score int) {
	if p.Table.DeskID != r.CallDeskID {
		panic(errors.New("you can't call score"))
	}
	p.Call(score)
	if r.MaxCallSocre < score {
		r.MaxCallSocre = score
	}
	r.CallScoreNum += 1
}

// 身份确认
func (r *Room) ConfirmRole() {
	for i, v := range r.Players {
		if v.Table.CallScore == r.MaxCallSocre {
			v.ConfirmRole(2)
			r.CallDeskID = i
		} else {
			v.ConfirmRole(1)
		}
	}
}

// 洗牌
func ShuffleDeck(deck []Card) {
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
}

// 出牌
func (r *Room) PlayCards(p *Player, cards []Card) {
	if r.BeforePlayDesk == r.CallDeskID {
		r.BeforeCards = r.BeforeCards[:0]
	}
	if !IsValidPlay(r.BeforeCards, cards) {
		panic(errors.New("can't play cards"))
	}
	if r.CallDeskID != p.Table.DeskID {
		panic(errors.New("not your turn"))
	}
	p.PlayCards(cards)
	if len(p.Cards) <= 0 {
		r.winRole = p.Table.Role
		r.GameOver()
	}
	if len(cards) > 0 {
		r.BeforePlayDesk = r.CallDeskID
	}
	r.BeforeCards = cards
	r.CallConvert()
}

// 结算
func (r *Room) Settle() {
	for _, p := range r.Players {
		p.Table.IsWin = r.winRole == p.Table.Role
		if p.Table.IsWin == false {
			var bet uint64
			if r.winRole == 1 {
				bet = utils.MaxUint64(p.Ctx.User.Coin, uint64(r.Score))
			} else {
				bet = utils.MaxUint64(p.Ctx.User.Coin, uint64(r.Score*2))
			}
			r.CoinPool += bet
			p.Ctx.User.Coin -= bet
		}
	}

	for _, p := range r.Players {
		if p.Table.IsWin == true {
			if p.Table.Role == 1 {
				p.Ctx.User.Coin += r.CoinPool / 2
			} else {
				p.Ctx.User.Coin += r.CoinPool
			}
		}
	}
}

// 游戏结束
func (r *Room) GameOver() {
	r.Settle()
	r.IsOver = true
	for _, p := range r.Players {
		p.Reset()
	}
}

// 关闭房间
func (r *Room) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.IsClosed = true
}

// 判断一个玩家是否在房间中
func (r *Room) InRoom(uid uint) bool {
	for _, v := range r.Players {
		if v != nil && v.Table.ID == uid {
			return true
		}
	}
	return false
}

func (r *Room) GetRet(p *Player) utils.Dict {
	ret := utils.Dict{
		"id":       r.ID,
		"is_over":  r.IsOver,
		"win_role": r.winRole,
		"score":    r.Score,
	}

	ret["myInfo"] = utils.Dict{
		"cards": p.Cards,
	}
	playersInfo := make([]utils.Dict, 0)
	for _, p := range r.Players {
		pInfo := utils.Dict{}
		if p != nil {
			pInfo["id"] = p.Table.ID
			pInfo["role"] = p.Table.Role
			pInfo["is_win"] = p.Table.IsWin
			pInfo["is_ready"] = p.Table.Ready
			pInfo["card_num"] = len(p.Cards)
			pInfo["coin"] = p.Ctx.User.Coin
			pInfo["name"] = p.Ctx.User.Name
		}
		playersInfo = append(playersInfo, pInfo)
	}
	ret["players"] = playersInfo
	return ret
}
