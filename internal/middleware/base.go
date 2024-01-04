package middleware

import (
	"my_app/internal/src"
	"my_app/internal/utils"
)

var MiddlewareList []Middleware

type Middleware interface {
	BeforeHandle(ctx *src.Ctx, data utils.Dict) utils.Dict
	AfterHandle(ctx *src.Ctx, ret utils.Dict) utils.Dict
}

func RegisterMiddleware(middleware Middleware) {
	MiddlewareList = append(MiddlewareList, middleware)
}

func init() {
	RegisterMiddleware(&LoginMiddleware{})
}
