package src

import (
	"sort"
)

type CardsType = int

const (
	Unknown           CardsType = iota // 未知牌型
	Single                             // 单张
	Pair                               // 对子
	Three                              // 三张
	Straight                           // 顺子
	PairStraight                       // 连对
	ThreeWithOne                       // 三带一
	ThreeWithTwo                       // 三带二
	Bomb                               // 炸弹
	KingBomb                           // 王炸
	PlaneWithoutWings                  // 飞机不带翅膀
	PlaneWithSingle                    // 飞机带单牌
	PlaneWithPair                      // 飞机带对子
)

var CardWeight = map[string]int{
	"3":          0,
	"4":          1,
	"5":          2,
	"6":          3,
	"7":          4,
	"8":          5,
	"9":          6,
	"10":         7,
	"J":          8,
	"Q":          9,
	"K":          10,
	"A":          11,
	"2":          12,
	"SmallJoker": 13,
	"BigJoker":   14,
}

// 定义牌的类型
type Card struct {
	Suit  string // 花色
	Value string // 点数
}

// 创建一副牌
func NewCards() []Card {
	suits := []string{"Spades", "Hearts", "Diamonds", "Clubs"}
	values := []string{"3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A", "2"}
	jokers := []string{"BigJoker", "SmallJoker"}

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

// 获取卡片的权重
func (card *Card) GetWeight() int {
	return CardWeight[card.Value]
}

// 判断是否为对子
func isPair(cards []Card) bool {
	if len(cards) != 2 {
		return false
	}
	return cards[0].GetWeight() == cards[1].GetWeight()
}

// 判断是否为三张
func isThree(cards []Card) bool {
	if len(cards) != 3 {
		return false
	}
	if cards[0].Value == cards[1].Value && cards[1].Value == cards[2].Value {
		return true
	}
	return false
}

// 判断是否为三带一
func isThreeWithOne(cards []Card) bool {
	if len(cards) != 4 {
		return false
	}

	sort.Slice(cards, func(i, j int) bool {
		return cards[i].GetWeight() < cards[j].GetWeight()
	})

	if cards[0].Value == cards[1].Value && cards[1].Value == cards[2].Value {
		return true
	}

	if cards[1].Value == cards[2].Value && cards[2].Value == cards[3].Value {
		return true
	}

	return false
}

// 判断是否为三带二
func isThreeWithTwo(cards []Card) bool {
	if len(cards) != 5 {
		return false
	}

	sort.Slice(cards, func(i, j int) bool {
		return cards[i].GetWeight() < cards[j].GetWeight()
	})

	if cards[0].Value == cards[1].Value && cards[1].Value == cards[2].Value && cards[3].Value == cards[4].Value {
		return true
	}

	if cards[0].Value == cards[1].Value && cards[2].Value == cards[3].Value && cards[3].Value == cards[4].Value {
		return true
	}

	return false
}

// 判断是否为顺子
func isStraight(cards []Card) bool {
	if len(cards) < 5 {
		return false
	}
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].GetWeight() < cards[j].GetWeight()
	})
	for i := 0; i < len(cards)-1; i++ {
		if cards[i].GetWeight() > 11 {
			return false
		}
		if cards[i].GetWeight()+1 != cards[i+1].GetWeight() {
			return false
		}
	}
	return true
}

// 判断是否为炸弹
func isBomb(cards []Card) bool {
	if len(cards) != 4 {
		return false
	}

	for i := 1; i < len(cards); i++ {
		if cards[i].Value != cards[0].Value {
			return false
		}
	}
	return true
}

// 判断是否为王炸
func isKingBomb(cards []Card) bool {
	if len(cards) != 2 {
		return false
	}
	for _, card := range cards {
		if card.Value != "Big Joker" && card.Value != "Small Joker" {
			return false
		}
	}
	return true
}

// 判断是否为连对
func isPairStraight(cards []Card) bool {
	if len(cards) < 6 || len(cards)%2 != 0 {
		return false
	}

	sort.Slice(cards, func(i, j int) bool {
		return cards[i].GetWeight() < cards[j].GetWeight()
	})

	for i := 0; i < len(cards)-1; i += 2 {
		if cards[i].GetWeight() != cards[i+1].GetWeight() {
			return false
		}
		if i > 0 && cards[i].GetWeight() != cards[i-1].GetWeight()+1 {
			return false
		}
	}

	return true
}

// 判读飞机的数量 最小飞机值
func PlaneInfo(cards []Card) (int, int) {
	countMap := make(map[Card]int)
	planes := make([]Card, 0)
	for _, card := range cards {
		countMap[card]++
		if countMap[card] == 3 && card.GetWeight() <= 11 {
			planes = append(planes, card)
		}
	}

	if len(planes) == 0 {
		return 0, 0
	}

	sort.Slice(planes, func(i, j int) bool {
		return planes[i].GetWeight() < planes[j].GetWeight()
	})

	minPlane := 0
	maxLength := 1
	currentLength := 1
	for i := 1; i < len(planes); i++ {
		if planes[i-1].GetWeight()+1 == planes[i].GetWeight() {
			currentLength++
		} else {
			currentLength = 1
		}
		if currentLength > maxLength {
			maxLength = currentLength
			minPlane = planes[i-currentLength+1].GetWeight()
		}
	}
	return maxLength, minPlane
}

// 判断是否为飞机不带翅膀
func isPlaneWithoutWings(cards []Card) bool {
	if len(cards) < 6 || len(cards)%3 != 0 {
		return false
	}
	planeNum, _ := PlaneInfo(cards)
	return planeNum == len(cards)%3
}

// 判断是否为飞机带单牌
func isPlaneWithSingle(cards []Card) bool {
	if len(cards)%4 != 0 || len(cards) < 8 {
		return false
	}
	planeNum, _ := PlaneInfo(cards)
	return planeNum >= len(cards)%4
}

// 判断是否为飞机带对子
func isPlaneWithPair(cards []Card) bool {
	if len(cards)%5 != 0 || len(cards) < 10 {
		return false
	}

	countMap := make(map[string]int)
	for _, card := range cards {
		countMap[card.Value]++
	}
	pairs := 0
	planes := 0
	for _, count := range countMap {
		if count == 2 {
			pairs++
		} else if count == 3 {
			planes++
		}
	}
	return pairs == planes && planes == len(cards)%5
}

// 判读出牌类型
func GetCardsType(cards []Card) CardsType {
	switch len(cards) {
	case 1: // 单张
		return Single
	case 2: // 对子或王炸
		if isPair(cards) {
			return Pair
		} else if isKingBomb(cards) {
			return KingBomb
		}
	case 3: // 三张
		if isThree(cards) {
			return Three
		}
	case 4: // 炸弹 三带一
		if isThreeWithOne(cards) {
			return ThreeWithOne
		} else if isBomb(cards) {
			return Bomb
		}
	default: // 顺子 连对 飞机 三带二 飞机
		if isThreeWithTwo(cards) {
			return ThreeWithTwo
		} else if isPairStraight(cards) {
			return PairStraight
		} else if isStraight(cards) {
			return Straight
		} else if isPlaneWithoutWings(cards) {
			return PlaneWithoutWings
		} else if isPlaneWithSingle(cards) {
			return PlaneWithSingle
		} else if isPlaneWithPair(cards) {
			return PlaneWithPair
		}
	}
	return Unknown
}

// 是否有效的出牌
func IsValidPlay(bCards []Card, cards []Card) bool {
	cardType := GetCardsType(cards)
	beforeCardType := GetCardsType(bCards)
	if len(cards) == 0 {
		return true
	}
	if cardType == Unknown {
		return false
	}
	if len(bCards) == 0 {
		return true
	}
	if cardType != beforeCardType && (cardType != Bomb || cardType != KingBomb) {
		return false
	}
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].GetWeight() < cards[j].GetWeight()
	})

	if cardType == beforeCardType {
		switch cardType {
		case Single, Pair, Three:
			return cards[0].GetWeight() > bCards[0].GetWeight()
		case Straight:
			return len(bCards) == len(cards) && cards[0].GetWeight() > bCards[0].GetWeight()
		case PairStraight:
			return len(bCards) == len(cards) && cards[0].GetWeight() > bCards[0].GetWeight()
		case PlaneWithPair, PlaneWithSingle, PlaneWithoutWings, ThreeWithOne, ThreeWithTwo:
			_, minPlane := PlaneInfo(cards)
			_, bMinPlane := PlaneInfo(bCards)
			return len(bCards) == len(cards) && minPlane > bMinPlane
		case Bomb:
			return cards[0].GetWeight() > bCards[0].GetWeight()
		}
	} else {
		switch cardType {
		case Bomb:
			return beforeCardType != KingBomb
		case KingBomb:
			return true
		}
	}
	return false
}
