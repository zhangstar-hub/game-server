package src

import (
	"fmt"
	"my_app/internal/ctx"
	"my_app/internal/utils"
)

func Test(ctx *ctx.Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	ret["test"] = "test"
	fmt.Printf("ctx: %#v\n", *ctx)
	fmt.Printf("user: %#v\n", (*ctx).User)
	return ret
}
