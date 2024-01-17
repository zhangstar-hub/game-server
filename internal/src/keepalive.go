package src

import (
	"my_app/internal/utils"
	"time"
)

// 心跳检测
func ReqKeepAlive(ctx *Ctx, data utils.Dict) utils.Dict {
	ret := make(utils.Dict)
	ctx.LastActiveTime = time.Now()
	return ret
}
