package doudizhu

// 房间
type Room struct {
	Cards   []Card   // 桌面上的纸牌
	Players []Player // 玩家
}

func NewRoom() *Room {
	return &Room{}
}
