package utils

import (
	"fmt"
	"runtime"
	"sync"
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

// 计算map长度
func MapLength(m map[string]interface{}) int {
	var size = 0
	for range m {
		size += 1
	}
	return size
}

// 计算SynMap长度
func SynMapLength(m *sync.Map) int {
	var size = 0

	m.Range(func(key, value interface{}) bool {
		size += 1
		return true // 返回 true 继续遍历
	})
	return size
}
