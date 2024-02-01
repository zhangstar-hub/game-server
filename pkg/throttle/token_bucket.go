package throttle

import (
	"context"
	"sync"
	"time"
)

type TokenBucketThrottle struct {
	mu            sync.Mutex
	tokenCount    uint32
	maxTokens     uint32
	tokenInterval time.Duration
	ctx           context.Context
	cancel        context.CancelFunc
}

// 创建一个令牌通
func NewTokenBucketThrottle() *TokenBucketThrottle {
	var maxTokens uint32 = 100
	t := &TokenBucketThrottle{
		maxTokens:     maxTokens,
		tokenInterval: 10 * time.Microsecond,
	}
	t.ctx, t.cancel = context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(t.tokenInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				t.mu.Lock()
				if t.tokenCount < t.maxTokens {
					t.tokenCount++
				}
				t.mu.Unlock()
			case <-t.ctx.Done():
				return
			}
		}
	}()
	return t
}

// 关闭令牌通限制
func (t *TokenBucketThrottle) Close() {
	t.cancel()
}

// 请求是否可以通过
func (t *TokenBucketThrottle) CanRequest() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.tokenCount > 0 {
		t.tokenCount--
		return true
	}
	return false
}
