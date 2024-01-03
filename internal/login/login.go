package login

import (
	"my_app/internal/models"
	"strings"
)

func Login(data map[string]interface{}) (ret map[string]interface{}) {
	ret = make(map[string]interface{})

	name := strings.TrimSpace(data["name"].(string))
	password := strings.TrimSpace(data["password"].(string))
	user := models.GetUserByName(name, password)
	if user == nil {
		user = models.CreateUser(name, password)
	}
	ret["user"] = map[string]interface{}{
		"id":   user.ID,
		"name": user.Name,
		"coin": user.Coin,
	}
	return ret
}
