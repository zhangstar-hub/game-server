package middleware

import (
	"errors"
	"my_app/internal/src"
	"my_app/internal/utils"
	"time"
)

type LoginMiddleware struct{}

func (m *LoginMiddleware) BeforeHandle(ctx *src.Ctx, data utils.Dict) utils.Dict {
	if ctx.Cmd != "login" && ctx.User == nil {
		panic(errors.New("login required"))
	}
	return data
}

func (m *LoginMiddleware) AfterHandle(ctx *src.Ctx, ret utils.Dict) utils.Dict {
	ctx.LastActiveTime = time.Now()
	return ret
}
