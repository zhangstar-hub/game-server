package main

import (
	"fmt"
	"my_app/internal/utils"
	"strings"
	"time"
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

func main() {
	fmt.Printf("utils.MondayFlushTime(): %v\n", utils.MondayFlushTime().Add(time.Second))
}
