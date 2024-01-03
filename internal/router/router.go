package router

import "my_app/internal/login"

type viewFunction func(data map[string]interface{}) map[string]interface{}

var Routers = map[string]viewFunction{
	"login": login.Login,
}
