package models

type PlayerModel struct {
	ID        uint   `gorm:"primary_key"`
	Cards     string `gorm:"type:json"` // 拥有的手牌
	RoomID    int    // 房间ID
	Role      int    // 角色 1：农民 2：地主
	Ready     bool   // 是否准备
	CallScore int    // 叫分
	IsWin     bool   // 是否胜利
	IsCall    bool   // 是否在出牌阶段
}
