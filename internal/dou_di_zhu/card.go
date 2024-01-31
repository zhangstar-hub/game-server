package doudizhu

import (
	"math/rand"
	"sort"
)

// 定义牌的类型
type Card struct {
	Suit  string // 花色
	Value string // 点数
}

// 创建一副牌
func NewDeck() []Card {
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

// 洗牌
func ShuffleDeck(deck []Card) []Card {
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	return deck
}

// 发牌
func DealCards(players []Player, deck []Card) {
	for i := 0; i < len(players); i++ {
		players[i].Cards = deck[i*17 : (i+1)*17]
		sort.Slice(players[i].Cards, func(j, k int) bool {
			return players[i].Cards[j].Value < players[i].Cards[k].Value
		})
	}
}
