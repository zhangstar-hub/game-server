package src

import (
	"fmt"
	"my_app/internal/models"
	"my_app/internal/utils"
	"strings"
	"time"
)

func ReqLogin(ctx *Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)

	name := strings.TrimSpace(data["name"].(string))
	password := strings.TrimSpace(data["password"].(string))
	user, err := models.GetUserByName(name, password)
	if err != nil {
		ret["error"] = err.Error()
		return ret
	}
	fmt.Printf("user: %#v\n", user)
	if user == nil {
		user = models.CreateUser(name, password)
	}
	fmt.Printf("user: %#v\n", user)
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
	ret["user"] = utils.Dict{
		"id":   user.ID,
		"name": user.Name,
		"coin": user.Coin,
	}
	return ret
}
