package zmq_client

import "my_app/internal/utils"

type viewFunction func(*ZMQClient, utils.Dict)

var ZMQRouters = map[string]viewFunction{
	"ReqZTest":        ReqZTest,
	"ReqZUserExit":    ReqZUserExit,
	"ReqZFlushConfig": ReqZFlushConfig,
	"ReqZRoomReady":   ReqZRoomReady,
	"ReqZPlayCards":   ReqZPlayCards,
	"ReqZEnterRoom":   ReqZEnterRoom,
	"ReqZCallScore":   ReqZCallScore,
}
