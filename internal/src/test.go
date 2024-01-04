package src

import (
	"fmt"
	"my_app/internal/utils"
)

func Test(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	ret["test"] = "test"
	fmt.Printf("ctx: %#v\n", *ctx)
	fmt.Printf("user: %#v\n", (*ctx).User)
	return ret
}

func ReqAddCoin(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	ctx.User.Coin += uint64(data["coin"].(float64))
	ret["coin"] = ctx.User.Coin
	return ret
}
