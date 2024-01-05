package src

import (
	"errors"
	"my_app/internal/context"
	"my_app/internal/models"
	"my_app/internal/utils"
	"my_app/internal/zmq_client"
	"strings"
	"time"
)

func ReqLogin(ctx *context.Ctx, data utils.Dict) (ret utils.Dict) {
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
	startTime := time.Now()
	for ok := ctx.IsOnline(user.ID); ok; {
		if time.Since(startTime) > 10*time.Second {
			panic(errors.New("login timeout"))
		}
		zmq_client.ZClient.Send(map[string]interface{}{
			"cmd": "ReqUserExit",
			"data": map[string]interface{}{
				"uid": user.ID,
			},
		})
		time.Sleep(200 * time.Millisecond)
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
