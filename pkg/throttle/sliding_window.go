package throttle

import (
	"time"
)

type SlidingWindowThrottle struct {
	windowSize  time.Duration // 时间窗口大小
	maxRequests int           // 时间窗口内最大请求数
	requests    []time.Time   // 请求时间队列
	closeFlag   bool          // 关闭信号
}

func NewSlidingWindowThrottle() *SlidingWindowThrottle {
	return &SlidingWindowThrottle{
		windowSize:  1 * time.Second,
		maxRequests: 100,
		requests:    make([]time.Time, 0),
		closeFlag:   false,
	}
}

// 关闭
func (t *SlidingWindowThrottle) Close() {
	t.closeFlag = true
	t.requests = t.requests[:0]
}

// 滑动窗口限制
func (t *SlidingWindowThrottle) CanRequest() bool {
	if t.closeFlag {
		return true
	}
	now := time.Now()
	for len(t.requests) > 0 && now.Sub(t.requests[0]) > t.windowSize {
		t.requests = t.requests[1:]
	}
	if len(t.requests) >= t.maxRequests {
		return false
	}
	t.requests = append(t.requests, now)
	return true
}
