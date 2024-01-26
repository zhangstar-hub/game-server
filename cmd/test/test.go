package main

import (
	"fmt"
	"reflect"
	"strings"
)

func AddString(a, b string) {
	_ = a + b
}

func SprintString(a, b string) {
	fmt.Sprintf("%s%s", a, b)
}

func BuildString(a, b string) {
	var builder strings.Builder
	builder.WriteString(a)
	builder.WriteString(b)
	_ = a + b
}

func interfaceIsNil(x interface{}) {
	if x == nil {
		fmt.Println("empty interface")
		return
	}
	fmt.Println("non-empty interface")
}

func main() {
	x := 1
	v := reflect.ValueOf(&x)
	fmt.Printf("v.Interface(): %v\n", v.Elem().CanSet())
	v.Elem().SetInt(2)
	fmt.Printf("x: %v\n", x)
}
