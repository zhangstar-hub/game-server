package throttle

import (
	"sync/atomic"
	"time"
)

type TokenBucketThrottle struct {
	tokenBucket   chan struct{} // 桶
	maxTokens     uint          // 最大令牌数量
	tokenInterval time.Duration // 令牌产生时间间隔
	closeFlag     atomic.Value  // 关闭
}

// 创建一个令牌通
func NewTokenBucketThrottle() *TokenBucketThrottle {
	t := &TokenBucketThrottle{
		tokenBucket:   make(chan struct{}),
		maxTokens:     1000,
		tokenInterval: 10 * time.Microsecond,
		closeFlag:     atomic.Value{},
	}
	t.closeFlag.Store(false)
	go func() {
		ticker := time.NewTicker(t.tokenInterval)
		defer ticker.Stop()
		for range ticker.C {
			if t.closeFlag.Load() == true {
				return
			}
			t.tokenBucket <- struct{}{}
		}
	}()
	return t
}

// 关闭令牌通限制
func (t *TokenBucketThrottle) Close() {
	close(t.tokenBucket)
	t.closeFlag.Store(true)
}

// 请求是否可以通过
func (t *TokenBucketThrottle) CanRequest() bool {
	<-t.tokenBucket
	return true
}
