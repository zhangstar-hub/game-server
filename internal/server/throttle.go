package server

import (
	"time"
)

const (
	MaxTokens     = 1000                  // 最大令牌数量
	TokenInterval = 10 * time.Millisecond // 令牌产生时间间隔

	windowSize  = 1 * time.Second // 时间窗口大小
	maxRequests = 100             // 时间窗口内最大请求数
)

var tokenBucket chan bool
var requests []time.Time

// 令牌桶
func TokenBucketLimit() {
	ticker := time.NewTicker(TokenInterval)
	tokenBucket = make(chan bool, MaxTokens)
	for range ticker.C {
		tokenBucket <- true
	}
}

// 滑动窗口限制
func SlidingWindowLimit() bool {

	now := time.Now()

	for len(requests) > 0 && now.Sub(requests[0]) > windowSize {
		requests = requests[1:]
	}

	if len(requests) >= maxRequests {
		return false
	}

	requests = append(requests, now)
	return true
}

// 请求阻塞
func CanRequest() {
	<-tokenBucket
	for !SlidingWindowLimit() {
		time.Sleep(100 * time.Millisecond)
	}
}

func init() {
	go TokenBucketLimit()
	go SlidingWindowLimit()
}
