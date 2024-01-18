package utils

// 查找字符串
func ArrayIndexOfString(data []string, target string) int {
	for i := 0; i < len(data); i++ {
		if target == data[i] {
			return i
		}
	}
	return -1
}
