package utils

import "sync"

// 计算map长度
func MapLength(m map[string]interface{}) int {
	var size = 0
	for range m {
		size += 1
	}
	return size
}

// 合并map
func MergeMaps(map1, map2 Dict) {
	for key, value := range map2 {
		map1[key] = value
	}
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
