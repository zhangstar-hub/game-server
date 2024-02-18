package src

import (
	"fmt"
	"math/rand"
	"my_app/internal/utils"
	"sort"
	"sync"
)

// 房间
type Room struct {
	ID       int32 // 房间ID
	Ctx      *Ctx
	Cards    []Card     // 桌面上的纸牌
	Players  []*Player  // 玩家
	Score    uint64     // 分数
	CoinPool uint64     //	奖池
	mu       sync.Mutex // 一个锁
	IsOver   bool       // 游戏是否结束
	IsFull   bool       // 是否满房了
	IsClosed bool       // 房间是否关闭
	winRole  int        // 胜利的角色
}

func NewRoom(ctx *Ctx) *Room {
	return &Room{
		Ctx:      ctx,
		Cards:    NewCards(),
		Players:  make([]*Player, 0),
		Score:    0,
		CoinPool: 0,
		mu:       sync.Mutex{},
		IsOver:   false,
		IsFull:   false,
		IsClosed: false,
		winRole:  0,
	}
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
	r.Players = append(r.Players, p)
	fmt.Printf("%v\n", r.Players)
	p.Table.RoomID = r.ID
	if len(r.Players) == 3 {
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
			r.Players = append(r.Players[:i], r.Players[i+1:]...)
			r.IsFull = false
			p.Reset()
			break
		}
	}
}

// 发牌
func (r *Room) DealCards() {
	for i := 0; i < len(r.Players); i++ {
		r.Players[i].Cards = append([]Card{}, r.Cards[i*17:(i+1)*17]...)
		sort.Slice(r.Players[i].Cards, func(j, k int) bool {
			return r.Players[i].Cards[j].Value < r.Players[i].Cards[k].Value
		})
	}
}

// 洗牌
func ShuffleDeck(deck []Card) {
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
}

// 出牌
func (r *Room) PlayCard(cards []Card) {
	for i, p := range r.Players {
		if p.Table.IsCall {
			if err := p.PlayCard(cards); err != nil {
				panic(err)
			}
			if len(p.Cards) <= 0 {
				r.winRole = p.Table.Role
				r.GameOver()
				return
			}
			p.Table.IsCall = false
			r.Players[(i+1)%3].Table.IsCall = true
			break
		}
	}
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
		if v.Table.ID == uid {
			return true
		}
	}
	return false
}

func (r *Room) GetRet() utils.Dict {
	ret := utils.Dict{
		"id":       r.ID,
		"is_over":  r.IsOver,
		"is_full":  r.IsFull,
		"win_role": r.winRole,
		"score":    r.Score,
	}

	ret["myInfo"] = utils.Dict{
		"cards": r.Ctx.Player.Cards,
	}
	playersInfo := make([]utils.Dict, 0)
	for _, p := range r.Players {
		playersInfo = append(playersInfo, utils.Dict{
			"id":       p.Table.ID,
			"role":     p.Table.Role,
			"is_call":  p.Table.IsCall,
			"is_win":   p.Table.IsWin,
			"is_ready": p.Table.Ready,
			"card_len": len(p.Cards),
			"coin":     p.Ctx.User.Coin,
			"name":     p.Ctx.User.Name,
		})
	}
	ret["players"] = playersInfo
	return ret
}
