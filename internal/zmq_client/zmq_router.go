package zmq_client

import "my_app/internal/utils"

type viewFunction func(utils.Dict)

var ZMQRouters = map[string]viewFunction{
	"ReqTest":     ReqTest,
	"ReqUserExit": ReqUserExit,
}
