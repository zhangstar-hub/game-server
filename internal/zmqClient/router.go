package zmqClient

type viewFunction func(string)

var Routers = map[string]viewFunction{}
