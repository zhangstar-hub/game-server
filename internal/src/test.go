package src

import (
	"my_app/internal/config"
	"my_app/internal/utils"
)

func ReqTest(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	ret["test"] = "test"
	ret["config"] = config.GetC()
	ret["params"] = utils.Dict{}
	utils.MergeMaps(ret["params"].(utils.Dict), data)
	return ret
}

func ReqAddCoin(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	AddCoin(ctx.User.ID, int64(data["coin"].(float64)))
	ret["coin"] = ctx.User.Coin
	return ret
}

func ReqZmqTest(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ctx.ZClient.Send(map[string]interface{}{"test": "test"})
	return ret
}
