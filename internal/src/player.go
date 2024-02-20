package src

import (
	"errors"
	"my_app/internal/utils"

	"github.com/thoas/go-funk"
)

type Player struct {
	ID        uint   // 玩家ID
	Cards     []Card // 拥有的手牌
	DeskID    int    // 座位号 0/1/2
	Role      int    // 角色 1：农民 2：地主
	Ready     bool   // 是否准备
	IsWin     bool   // 是否胜利
	CallScore int    // 叫分
	Coin      uint64 //玩家的金币数量
	Name      uint64 //玩家的名字
}

// 获取玩家数据
func PlayerLoadData(ctx *Ctx) {
	var player *Player
	if ctx.User.RoomID != 0 {
		if ok, r := ctx.RoomManager.GetRoom(ctx); ok {
			for _, p := range r.Players {
				if p.ID == ctx.User.ID {
					continue
				}
				player = p
				break
			}
		}
	} else {
		player = &Player{
			ID:        ctx.User.ID,
			Cards:     []Card{},
			DeskID:    0,
			Role:      1,
			Ready:     false,
			CallScore: 0,
		}
	}
	ctx.Player = player
}

// 准备
func (p *Player) SetReady(status bool) {
	p.Ready = status
}

// 叫分
func (p *Player) Call(score int) {
	p.CallScore = score
}

// 身份确认
func (p *Player) ConfirmRole(role int) {
	p.Role = role
}

// 打牌
func (p *Player) PlayCards(cards []Card) {
	for _, card := range cards {
		if !funk.Contains(p.Cards, card) {
			panic(errors.New("play card error"))
		}
	}
	p.Cards = funk.Join(p.Cards, cards, funk.LeftJoin).([]Card)
}

// 重置对局
func (p *Player) Reset() {
	p.Cards = []Card{}
	p.Ready = false
	p.Role = 1
	p.CallScore = 0
}

// 数据获取
func (p *Player) GetRet() (ret utils.Dict) {
	ret = utils.Dict{
		"cards": p.Cards,
		"role":  p.Role,
		"score": p.CallScore,
		"ready": p.Ready,
	}
	return ret
}
