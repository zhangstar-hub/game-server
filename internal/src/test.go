package src

import (
	"my_app/internal/context"
	"my_app/internal/utils"
	"my_app/internal/zmq_client"
)

func ReqTest(ctx *context.Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	ret["test"] = "test"
	return ret
}

func ReqAddCoin(ctx *context.Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	ctx.User.Coin += uint64(data["coin"].(float64))
	ret["coin"] = ctx.User.Coin
	return ret
}

func ReqZmqTest(ctx *context.Ctx, data utils.Dict) (ret utils.Dict) {
	zmq_client.ZClient.Send(map[string]interface{}{"test": "test"})
	return ret
}
