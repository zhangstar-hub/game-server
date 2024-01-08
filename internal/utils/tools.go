package utils

import (
	"fmt"
	"net"
	"os"
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

// 获取本机IP地址
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// 获取服务器进程ID
func GetPid() int {
	return int(os.Getpid())
}
