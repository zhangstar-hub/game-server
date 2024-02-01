package src

// 定义牌的类型
type Card struct {
	Suit  string // 花色
	Value string // 点数
}

// 创建一副牌
func NewCards() []Card {
	suits := []string{"Spades", "Hearts", "Diamonds", "Clubs"}
	values := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	jokers := []string{"Big Joker", "Small Joker"}

	deck := make([]Card, 0)

	// 添加普通牌
	for _, suit := range suits {
		for _, value := range values {
			card := Card{Suit: suit, Value: value}
			deck = append(deck, card)
		}
	}

	// 添加大小王
	for _, joker := range jokers {
		card := Card{Suit: "Joker", Value: joker}
		deck = append(deck, card)
	}
	return deck
}
