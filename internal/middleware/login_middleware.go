package middleware

import (
	"errors"
	"my_app/internal/src"
	"my_app/internal/utils"
	"time"

	"github.com/thoas/go-funk"
)

type LoginMiddleware struct{}

// 不需要登录就可以接受的请求
var NoLoingReqList = []string{
	"ReqLogin", "ReqKeepAlive",
}

func (m *LoginMiddleware) BeforeHandle(ctx *src.Ctx, data utils.Dict) utils.Dict {
	index := funk.IndexOfString(NoLoingReqList, ctx.Cmd)
	if index == -1 && ctx.User == nil {
		panic(errors.New("login required"))
	}
	return data
}

func (m *LoginMiddleware) AfterHandle(ctx *src.Ctx, ret utils.Dict) utils.Dict {
	ctx.LastActiveTime = time.Now()
	return ret
}
