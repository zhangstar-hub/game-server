package main

import "fmt"

const (
	Unknown           = iota // 未知牌型
	Single                   // 单张
	Pair                     // 对子
	Three                    // 三张
	Straight                 // 顺子
	PairStraight             // 连对
	ThreeWithOne             // 三带一
	ThreeWithTwo             // 三带二
	Bomb                     // 炸弹
	KingBomb                 // 王炸
	PlaneWithoutWings        // 飞机不带翅膀
	PlaneWithSingle          // 飞机带单牌
	PlaneWithPair            // 飞机带对子
)

func main() {
	a := interface{}(1)
	switch a {
	case '1', '2', interface{}(1):
		print("case 1:\n")
	case 1:
		print("case 2:\n")
	}
	b := []string{"1", "2", "3", "4", "5"}
	for _, i := range b {
		fmt.Printf("i: %v\n", i)
	}
}
