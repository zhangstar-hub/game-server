package router

import (
	"my_app/internal/ctx"
	"my_app/internal/src"
	"my_app/internal/utils"
)

type viewFunction func(ctx *ctx.Ctx, data utils.Dict) utils.Dict

var Routers = map[string]viewFunction{
	"login": src.Login,
	"test":  src.Test,
}
