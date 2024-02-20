package src

import (
	"encoding/json"
	"errors"
	"my_app/internal/db"
	"my_app/internal/models"
	"my_app/internal/utils"

	"github.com/thoas/go-funk"
)

type Player struct {
	Table *models.PlayerModel
	Cards []Card
	Ctx   *Ctx
}

// 初始数据
func PlayerInitData() models.PlayerModel {
	return models.PlayerModel{
		Cards:     "[]",
		Ready:     false,
		Role:      1,
		RoomID:    0,
		DeskID:    0,
		CallScore: 0,
	}
}

// 加载数据
func PlayerLoadData(ctx *Ctx) {
	Table := &models.PlayerModel{
		ID: ctx.User.ID,
	}
	db.DB.Attrs(PlayerInitData()).FirstOrInit(Table)
	ctx.Player = &Player{
		Table: Table,
		Ctx:   ctx,
		Cards: []Card{},
	}
	json.Unmarshal([]byte(ctx.Player.Table.Cards), ctx.Player.Cards)
}

// 保存数据
func (p *Player) Save() error {
	return db.DB.Save(p.Table).Error
}

// 准备
func (p *Player) SetReady(status bool) {
	p.Table.Ready = status
}

// 叫分
func (p *Player) Call(score int) {
	p.Table.CallScore = score
}

// 身份确认
func (p *Player) ConfirmRole(role int) {
	p.Table.Role = role
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
	p.Table.Ready = false
	p.Table.Role = 1
	p.Table.CallScore = 0
	p.Table.RoomID = 0
	p.Table.Cards = "[]"
}

// 数据获取
func (p *Player) GetRet() (ret utils.Dict) {
	ret = utils.Dict{
		"cards": p.Cards,
		"role":  p.Table.Role,
		"score": p.Table.CallScore,
		"ready": p.Table.Ready,
		"room":  p.Table.RoomID,
	}
	return ret
}
