package main

import (
	"fmt"
	"sort"
)

func isPairStraight(cards []int) bool {
	if len(cards) < 6 || len(cards)%2 != 0 {
		return false
	}

	sort.Slice(cards, func(i, j int) bool {
		return cards[i] < cards[j]
	})

	for i := 0; i < len(cards); i += 2 {
		if cards[i] != cards[i+1] {
			return false
		}
		if i+2 < len(cards)-1 && cards[i]+1 != cards[i+2] {
			return false
		}
	}

	return true
}

func main() {
	a := isPairStraight([]int{1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 7, 7})
	fmt.Printf("a: %v\n", a)
}
