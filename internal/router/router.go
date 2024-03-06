package router

import (
	"my_app/internal/src"
	"my_app/internal/utils"
)

type viewFunction func(ctx *src.Ctx, data utils.Dict) utils.Dict

var Routers = map[string]viewFunction{
	"ReqLogin":        src.ReqLogin,
	"ReqTest":         src.ReqTest,
	"ReqAddCoin":      src.ReqAddCoin,
	"ReqGetMission":   src.ReqGetMission,
	"ReqKeepAlive":    src.ReqKeepAlive,
	"ReqEnterRoom":    src.ReqEnterRoom,
	"ReqLeaveRoom":    src.ReqLeaveRoom,
	"ReqRoomReady":    src.ReqRoomReady,
	"ReqWatchCards":   src.ReqWatchCards,
	"ReqGetRole":      src.ReqGetRole,
	"ReqCallScore":    src.ReqCallScore,
	"ReqPlayCards":    src.ReqPlayCards,
	"ReqEnterNewRoom": src.ReqEnterNewRoom,
}
