package middleware

import (
	"errors"
	"my_app/internal/ctx"
	"my_app/internal/utils"
)

type LoginMiddleware struct{}

func (m *LoginMiddleware) BeforeHandle(ctx *ctx.Ctx, data utils.Dict) utils.Dict {
	if ctx.Cmd != "login" && ctx.User == nil {
		panic(errors.New("login required"))
	}
	return data
}

func (m *LoginMiddleware) AfterHandle(ctx *ctx.Ctx, ret utils.Dict) utils.Dict {
	return ret
}
