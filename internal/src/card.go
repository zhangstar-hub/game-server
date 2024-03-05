package src

import (
	"sort"
)

type CardsType = int

const (
	Unknown           CardsType = 0  // 未知牌型
	Single                      = 1  // 单张
	Pair                        = 2  // 对子
	Three                       = 3  // 三张
	Straight                    = 4  // 顺子
	PairStraight                = 5  // 连对
	ThreeWithOne                = 6  // 三带一
	ThreeWithTwo                = 7  // 三带二
	Bomb                        = 8  // 炸弹
	KingBomb                    = 9  // 王炸
	PlaneWithoutWings           = 10 // 飞机不带翅膀
	PlaneWithSingle             = 11 // 飞机带单牌
	PlaneWithPair               = 12 // 飞机带对子
	FourWithTow                 = 13 // 四带二
	FourWithTowPair             = 14 // 四带两对
)

// 定义牌的类型
type Card struct {
	Suit  int // 花色
	Value int // 点数
}

// 创建一副牌
func NewCards() []Card {
	suits := []int{1, 2, 3, 4}                                   // 方片 梅花 红桃 黑桃
	values := []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15} // 3 - K

	deck := make([]Card, 0)
	for _, suit := range suits {
		for _, value := range values {
			card := Card{Suit: suit, Value: value}
			deck = append(deck, card)
		}
	}
	deck = append(deck, Card{Suit: 5, Value: 16}, Card{Suit: 5, Value: 17})
	return deck
}

// 获取卡片值
func CardsToValue(cards []Card) []int {
	cardsValue := make([]int, len(cards))
	for i, v := range cards {
		cardsValue[i] = v.Suit*100 + v.Value
	}
	return cardsValue
}

// 卡片值转卡片对象
func ValueToCards(cardsVal []int) []Card {
	cards := make([]Card, len(cardsVal))
	for i, v := range cardsVal {
		cards[i] = Card{Suit: v / 100, Value: v % 100}
	}
	return cards
}

// 统计卡片类型数量
func countCardsType(cards []Card) [][2]int {
	cardsVal := make(map[int]int, 0)
	for _, card := range cards {
		cardsVal[card.Value] += 1
	}
	cardsNum := make([][2]int, 0, len(cardsVal))
	for key, count := range cardsVal {
		cardsNum = append(cardsNum, [2]int{key, count})
	}
	sort.Slice(cardsNum, func(i, j int) bool {
		if cardsNum[i][1] < cardsNum[j][1] {
			return true
		}
		if cardsNum[i][1] > cardsNum[j][1] {
			return false
		}
		return cardsNum[i][0] < cardsNum[j][0]
	})
	return cardsNum
}

// 判断是否为对子
func isPair(cards []Card) bool {
	if len(cards) != 2 {
		return false
	}
	return cards[0].Value == cards[1].Value
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
	cardsNum := countCardsType(cards)
	if len(cardsNum) != 2 {
		return false
	}
	return cardsNum[0][1] == 1
}

// 判断是否为三带二
func isThreeWithTwo(cards []Card) bool {
	if len(cards) != 5 {
		return false
	}
	cardsNum := countCardsType(cards)
	if len(cardsNum) != 2 {
		return false
	}
	return cardsNum[0][1] == 2
}

// 判断是否为顺子
func isStraight(cards []Card) bool {
	if len(cards) < 5 {
		return false
	}
	for i := 0; i < len(cards)-1; i++ {
		if cards[i].Value > 14 {
			return false
		}
		if cards[i].Value+1 != cards[i+1].Value {
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
		if card.Value != 16 && card.Value != 17 {
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

	cardsNum := countCardsType(cards)
	for i := 0; i < len(cardsNum); i++ {
		if cardsNum[i][0] > 14 {
			return false
		}
		if cardsNum[i][1] != 2 {
			return false
		}
		if i < len(cardsNum)-1 && cardsNum[i][0]+1 != cardsNum[i+1][0] {
			return false
		}
	}
	return true
}

// 判读3带的数量 最小飞机值
func PlaneInfo(cards []Card) (int, int) {
	cardsNum := countCardsType(cards)
	cardsNumWith3 := []int{}
	for _, cardNum := range cardsNum {
		if cardNum[1] >= 3 {
			cardsNumWith3 = append(cardsNumWith3, cardNum[0])
		}
	}
	if len(cardsNumWith3) <= 0 {
		return 0, 0
	}
	sort.Slice(cardsNumWith3, func(i, j int) bool {
		return cardsNumWith3[i] < cardsNumWith3[0]
	})

	maxPlaneNum := 0
	minPlaneVal := 0
	curPlaneNum := 1
	curMinPlaneVal := cardsNumWith3[0]

	for i := 1; i < len(cardsNumWith3); i++ {
		if cardsNumWith3[i] == cardsNumWith3[i-1]+1 {
			curPlaneNum += 1
			if curMinPlaneVal == 0 {
				curMinPlaneVal = cardsNumWith3[i]
			}
		} else {
			if curPlaneNum > maxPlaneNum {
				maxPlaneNum = curPlaneNum
				minPlaneVal = curMinPlaneVal
			}
			curPlaneNum = 1
			curMinPlaneVal = 0
		}
	}
	if curPlaneNum > maxPlaneNum {
		maxPlaneNum = curPlaneNum
		minPlaneVal = curMinPlaneVal
	}
	return maxPlaneNum, minPlaneVal
}

// 判断是否为飞机不带翅膀
func isPlaneWithoutWings(cards []Card) bool {
	if len(cards) < 6 || len(cards)%3 != 0 {
		return false
	}
	planeNum, _ := PlaneInfo(cards)
	return planeNum == len(cards)/3
}

// 判断是否为飞机带单牌
func isPlaneWithSingle(cards []Card) bool {
	if len(cards)%4 != 0 || len(cards) < 8 {
		return false
	}
	planeNum, _ := PlaneInfo(cards)
	return planeNum >= len(cards)/4
}

// 判断是否为飞机带对子
func isPlaneWithPair(cards []Card) bool {
	if len(cards)%5 != 0 || len(cards) < 10 {
		return false
	}

	countMap := make(map[int]int)
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
	return pairs == planes && planes == len(cards)/5
}

// 四带二
func isFourWithTow(cards []Card) bool {
	if len(cards) != 6 {
		return false
	}
	cardsNum := countCardsType(cards)
	return cardsNum[len(cards)-1][1] == 4
}

// 四带二对子
func isFourWithTowPair(cards []Card) bool {
	if len(cards) != 8 {
		return false
	}
	cardsNum := countCardsType(cards)
	if len(cardsNum) != 3 {
		return false
	}
	return cardsNum[0][1] == cardsNum[1][1] && cardsNum[2][1] == 4
}

// 判读出牌类型
func GetCardsType(cards []Card) CardsType {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Value < cards[j].Value
	})
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
	default: // 顺子 连对 飞机 三带二 飞机 四带二
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
		} else if isFourWithTow(cards) {
			return FourWithTow
		} else if isFourWithTowPair(cards) {
			return FourWithTowPair
		}
	}
	return Unknown
}

// 是否有效的出牌
func IsValidPlay(bCards []Card, cards []Card) bool {
	cardType := GetCardsType(cards)
	beforeCardType := GetCardsType(bCards)
	if cardType != beforeCardType && !(cardType == Bomb || cardType == KingBomb) {
		return false
	}
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Value < cards[j].Value
	})

	if cardType != beforeCardType {
		switch cardType {
		case Bomb:
			return beforeCardType != KingBomb
		case KingBomb:
			return true
		}
		return false
	}

	switch cardType {
	case Single, Pair, Three:
		return cards[0].Value > bCards[0].Value
	case Straight:
		return len(bCards) == len(cards) && cards[0].Value > bCards[0].Value
	case PairStraight:
		return len(bCards) == len(cards) && cards[0].Value > bCards[0].Value
	case PlaneWithPair, PlaneWithSingle, PlaneWithoutWings, ThreeWithOne, ThreeWithTwo:
		_, minPlane := PlaneInfo(cards)
		_, bMinPlane := PlaneInfo(bCards)
		return len(bCards) == len(cards) && minPlane > bMinPlane
	case Bomb:
		return cards[0].Value > bCards[0].Value
	}
	return false
}
