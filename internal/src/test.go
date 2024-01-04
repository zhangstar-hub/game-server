package src

import (
	"my_app/internal/utils"
)

func ReqTest(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	ret["test"] = "test"
	return ret
}

func ReqAddCoin(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	ctx.User.Coin += uint64(data["coin"].(float64))
	ret["coin"] = ctx.User.Coin
	return ret
}
