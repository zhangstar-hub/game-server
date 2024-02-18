package main

import "fmt"

type S struct {
	c  string
	c1 []int
}

func main() {
	var s *S
	if s == nil {
		fmt.Println("nil")
	}
}
