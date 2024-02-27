package main

import "fmt"

func main() {
	slice1 := make([]int, 10)
	fmt.Printf("slice1: %v\n", slice1)
	fmt.Printf("slice1: %p\n", slice1)

	slice1 = append(slice1, 4, 5)
	fmt.Printf("slice1: %v\n", slice1)
	fmt.Printf("slice1: %p\n", slice1)
}
