package router

import (
	"my_app/internal/context"
	"my_app/internal/src"
	"my_app/internal/utils"
)

type viewFunction func(ctx *context.Ctx, data utils.Dict) utils.Dict

var Routers = map[string]viewFunction{
	"ReqLogin":      src.ReqLogin,
	"ReqTest":       src.ReqTest,
	"ReqAddCoin":    src.ReqAddCoin,
	"ReqGetMission": src.ReqGetMission,
}
