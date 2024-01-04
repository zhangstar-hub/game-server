package src

import (
	"my_app/internal/models"
	"my_app/internal/utils"
	"strings"
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
	if user == nil {
		user = models.CreateUser(name, password)
	}
	ret["user"] = utils.Dict{
		"id":   user.ID,
		"name": user.Name,
		"coin": user.Coin,
	}
	ctx.User = user
	Users.Store(ctx.User.ID, ctx)
	return ret
}
