package middleware

import (
	"my_app/internal/context"
	"my_app/internal/utils"
)

var MiddlewareList []Middleware

type Middleware interface {
	BeforeHandle(ctx *context.Ctx, data utils.Dict) utils.Dict
	AfterHandle(ctx *context.Ctx, ret utils.Dict) utils.Dict
}

func RegisterMiddleware(middleware Middleware) {
	MiddlewareList = append(MiddlewareList, middleware)
}

func init() {
	RegisterMiddleware(&LoginMiddleware{})
}
