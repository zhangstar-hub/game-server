package main

import "fmt"

type S struct {
	c  string
	c1 []int
}

func main() {
	a := map[interface{}]struct{}{}

	s1 := S{c: "1", c1: []int{1}}
	s2 := S{c: "1", c1: []int{1}}
	s3 := S{c: "2", c1: []int{1}}
	a[s1] = struct{}{}
	a[s2] = struct{}{}
	a[s3] = struct{}{}
	fmt.Printf("a: %v\n", a)
}
