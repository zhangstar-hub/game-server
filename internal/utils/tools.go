package utils

import (
	"fmt"
	"runtime"
)

func PrintCallStack() {
	// 设置跟踪的深度，这里选择一个足够大的值以确保获取完整的调用栈信息
	const depth = 10
	// 创建一个足够大的切片来存储调用栈信息
	pc := make([]uintptr, depth)

	// 获取调用栈信息
	n := runtime.Callers(0, pc)
	if n == 0 {
		// 如果未获取到调用栈信息，则直接返回
		return
	}

	// 解析并打印调用栈信息
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		fmt.Printf("%s:%d %s\n", frame.File, frame.Line, frame.Function)
		if !more {
			break
		}
	}
}
