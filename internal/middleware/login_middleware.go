package middleware

import (
	"errors"
	"fmt"
	"my_app/internal/context"
	"my_app/internal/utils"
	"time"
)

type LoginMiddleware struct{}

// 不需要登录就可以接受的请求
var NoLoingReqList = []string{
	"ReqLogin",
}

func (m *LoginMiddleware) BeforeHandle(ctx *context.Ctx, data utils.Dict) utils.Dict {
	index := utils.ArrayIndexOfString(NoLoingReqList, ctx.Cmd)
	fmt.Printf("index: %v\n", index)
	if index == -1 && ctx.User == nil {
		panic(errors.New("login required"))
	}
	return data
}

func (m *LoginMiddleware) AfterHandle(ctx *context.Ctx, ret utils.Dict) utils.Dict {
	ctx.LastActiveTime = time.Now()
	return ret
}
