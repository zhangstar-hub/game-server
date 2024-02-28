package zmq_client

import (
	"encoding/json"
	"fmt"
	"my_app/internal/config"
	"my_app/internal/logger"
	"my_app/internal/src"
	"my_app/internal/utils"
	"sync"

	"github.com/pebbe/zmq4"
)

type ZMQClient struct {
	client      *zmq4.Socket
	CtxMap      *sync.Map
	RoomManager *src.RoomManager
}

func NewZMQClient(ctxMap *sync.Map, roomManager *src.RoomManager) *ZMQClient {
	conf := config.GetC()
	client, _ := zmq4.NewSocket(zmq4.DEALER)
	client.SetIdentity(fmt.Sprintf("%s:%d", utils.GetLocalIP(), conf.Env.App.Port))

	zClient := &ZMQClient{
		client:      client,
		CtxMap:      ctxMap,
		RoomManager: roomManager,
	}
	err := client.Connect(conf.Env.ZMQCenter.Address)
	if err != nil {
		panic(fmt.Sprintf("无法连接中心服务器：%s", conf.Env.ZMQCenter.Address))
	}
	_, err = client.Send("first_message", 0)
	if err != nil {
		panic("向中心服务器发送消息失败")
	}
	return zClient
}

// 向中心服务器发送数据
func (z *ZMQClient) Send(data utils.Dict) (int, error) {
	json, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	return z.client.SendBytes(json, 0)
}

// 消息定向广播
func (z *ZMQClient) BroastMessage(cmd string, from_uid uint, to_uid_list []uint, message utils.Dict) {
	z.Send(utils.Dict{
		"cmd": cmd,
		"data": utils.Dict{
			"from_uid":    from_uid,
			"to_uid_list": to_uid_list,
			"message":     message,
		},
	})
}

// 从中心服务器接受数据
func (z *ZMQClient) Recv() ([]byte, error) {
	return z.client.RecvBytes(0)
}

// 数据接口监听
func (z *ZMQClient) MessageListener() {
	for {
		message, err := z.Recv()
		fmt.Printf("MessageListener message: %s\n", message)
		if err != nil {
			fmt.Printf("message recv error %s\n", err)
			continue
		}
		messageMap := make(map[string]interface{})
		err = json.Unmarshal(message, &messageMap)
		if err != nil {
			fmt.Printf("message unmarshal error %s\n", err)
			continue
		}
		if _, ok := messageMap["cmd"]; !ok {
			fmt.Printf("message unmarshal error, cmd not found\n")
			continue
		}
		if _, ok := messageMap["data"]; !ok {
			fmt.Printf("message unmarshal error, data not found\n")
			continue
		}
		cmd := messageMap["cmd"].(string)
		if _, ok := ZMQRouters[cmd]; !ok {
			fmt.Printf("zmq cmd not found %s\n", cmd)
			continue
		}
		data := messageMap["data"].(utils.Dict)
		func() {
			defer func() {
				if err := recover(); err != nil {
					logger.ZMQInfo(fmt.Sprintf("ZMQClient MessageListener panic: %s", err))
				}
			}()
			ZMQRouters[cmd](z, data)
		}()
		msg := fmt.Sprintf("ClientRecv message:%s", message)
		logger.ZMQInfo(msg)
	}
}
