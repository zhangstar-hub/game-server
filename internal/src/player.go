package src

import (
	"errors"
	"fmt"
	"my_app/internal/utils"

	"github.com/thoas/go-funk"
)

type Player struct {
	ID        uint   // 玩家ID
	Cards     []Card // 拥有的手牌
	DeskID    int    // 座位号 0/1/2
	Role      int    // 角色 1：农民 2：地主
	Ready     bool   // 是否准备
	CallScore int    // 叫分
	Name      string // 玩家的名字
	Avatar    string // 头像
}

// 获取玩家数据
func PlayerLoadData(ctx *Ctx) {
	var player *Player
	if ctx.User.RoomID != 0 {
		if ok, r := ctx.RoomManager.GetRoom(ctx.User.RoomID, ctx.User.ID); ok {
			for _, p := range r.Players {
				if p.ID != ctx.User.ID {
					continue
				}
				player = p
				break
			}
		}
	}
	if player == nil {
		player = &Player{
			ID:        ctx.User.ID,
			Cards:     []Card{},
			DeskID:    0,
			Role:      1,
			Ready:     false,
			CallScore: -1,
			Name:      ctx.User.Name,
			Avatar:    ctx.User.Avatar,
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
	fmt.Printf("cards: %v\n", cards)
	fmt.Printf("p.Cards: %v\n", p.Cards)

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
	p.CallScore = -1
}

// 数据获取
func (p *Player) GetRet() (ret utils.Dict) {
	ret = utils.Dict{
		"id":         p.ID,
		"card_num":   len(p.Cards),
		"role":       p.Role,
		"call_score": p.CallScore,
		"is_ready":   p.Ready,
		"coin":       GetCoin(p.ID),
		"name":       p.Name,
		"avatar":     p.Avatar,
		"desk_id":    p.DeskID,
	}
	return ret
}
