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
func (p *Player) Ready() {
	p.Table.Ready = true
}

// 叫分
func (p *Player) Call(score int) {
	p.Table.CallScore = score
}

// 身份确认
func (p *Player) Confirm(role int) {
	p.Table.Role = role
}

// 打牌
func (p *Player) PlayCard(cards []Card) error {
	if funk.Every(p.Cards, cards) == false {
		return errors.New("play card error")
	}
	p.Cards = funk.Join(p.Cards, cards, funk.LeftJoin).([]Card)
	return nil
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
	ret = make(utils.Dict, 0)
	ret["cards"] = utils.Dict{
		"cards": p.Cards,
		"role":  p.Table.Role,
		"score": p.Table.CallScore,
		"ready": p.Table.Ready,
		"room":  p.Table.RoomID,
	}
	return
}
