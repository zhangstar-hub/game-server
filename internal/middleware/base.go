package middleware

import (
	"my_app/internal/ctx"
	"my_app/internal/utils"
)

var MiddlewareList []Middleware

type Middleware interface {
	BeforeHandle(ctx *ctx.Ctx, data utils.Dict) utils.Dict
	AfterHandle(ctx *ctx.Ctx, ret utils.Dict) utils.Dict
}

func RegisterMiddleware(middleware Middleware) {
	MiddlewareList = append(MiddlewareList, middleware)
}

func init() {
	RegisterMiddleware(&LoginMiddleware{})
}
