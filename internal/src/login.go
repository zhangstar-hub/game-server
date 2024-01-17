package src

import (
	"fmt"
	"my_app/internal/models"
	"my_app/internal/utils"
	"strings"
	"time"
)

// 登录接口
func ReqLogin(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)

	name := strings.TrimSpace(data["name"].(string))
	password := strings.TrimSpace(data["password"].(string))
	user, err := models.GetUserByName(name, password)
	if err != nil {
		ret["error"] = err.Error()
		return ret
	}
	if user == nil {
		user = models.CreateUser(name, password)
	}
	if ok := ctx.IsOnline(user.ID); ok {
		startTime := time.Now()
		ctx.QuitMessage(user.ID)
		for ok := ctx.IsOnline(user.ID); ok; ok = ctx.IsOnline(user.ID) {
			if time.Since(startTime) > 10*time.Second {
				fmt.Printf("force login: %d", user.ID)
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		user, _ = models.GetUserByName(name, password)
	}
	ctx.SetOnline(user.ID)
	ctx.User = user

	LoginLoadData(ctx)
	LoginCheckData(ctx, ret)
	LoginRetData(ctx, ret)

	ret["user"] = utils.Dict{
		"id":   user.ID,
		"name": user.Name,
		"coin": user.Coin,
	}
	return ret
}

// 登录加载数据
func LoginLoadData(ctx *Ctx) {
	LoginBonusLoadData(ctx)
}

// 登录检测数据
func LoginCheckData(ctx *Ctx, ret utils.Dict) {
	ctx.LoginBonusCtx.LoginCheck(ctx, ret)
}

// 登录返回数据
func LoginRetData(ctx *Ctx, ret utils.Dict) {
	ctx.LoginBonusCtx.GetRet(ret)
}
