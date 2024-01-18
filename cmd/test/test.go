package main

import (
	"fmt"
	"my_app/internal/utils"
	"time"
)

func main() {
	fmt.Printf("utils.MondayFlushTime(): %v\n", utils.MondayFlushTime().Add(time.Second))
}
