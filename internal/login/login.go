package login

import (
	"my_app/internal/models"
	"my_app/internal/utils"
	"strings"
)

func Login(data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)

	name := strings.TrimSpace(data["name"].(string))
	password := strings.TrimSpace(data["password"].(string))
	user := models.GetUserByName(name, password)
	if user == nil {
		user = models.CreateUser(name, password)
	}
	ret["user"] = utils.Dict{
		"id":   user.ID,
		"name": user.Name,
		"coin": user.Coin,
	}
	return ret
}
