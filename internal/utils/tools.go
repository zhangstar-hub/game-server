package utils

import (
	"fmt"
	"runtime"
)

// 查找字符串
func ArrayIndexOfString(data []string, target string) int {
	for i := 0; i < len(data); i++ {
		if target == data[i] {
			return i
		}
	}
	return -1
}

// 调用栈打印
func PrintStackTrace() {
	const size = 64 << 10 // 64KB
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, false)]

	fmt.Println(string(buf))
}
