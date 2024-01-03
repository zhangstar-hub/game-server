package router

import (
	"my_app/internal/login"
	"my_app/internal/utils"
)

type viewFunction func(data utils.Dict) utils.Dict

var Routers = map[string]viewFunction{
	"login": login.Login,
}
