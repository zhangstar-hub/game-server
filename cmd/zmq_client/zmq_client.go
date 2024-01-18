package main

import (
	"encoding/json"
	"fmt"
	"my_app/internal/utils"

	"github.com/pebbe/zmq4"
)

const (
	ZmqServerAddr = "tcp://127.0.0.1:5555"
)

type ZMQClient struct {
	client *zmq4.Socket
}

func NewZMQClient() *ZMQClient {
	client, _ := zmq4.NewSocket(zmq4.DEALER)
	client.SetIdentity(fmt.Sprintf("%s:%d", utils.GetLocalIP(), 6666))

	zClient := &ZMQClient{client: client}
	err := client.Connect(ZmqServerAddr)
	if err != nil {
		panic(fmt.Sprintf("无法连接中心服务器：%s", ZmqServerAddr))
	}
	_, err = client.Send("first_message", 0)
	if err != nil {
		panic("向中心服务器发送消息失败")
	}
	return zClient
}

// 向中心服务器发送数据
func (z *ZMQClient) Send(data map[string]interface{}) (int, error) {
	json, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	return z.client.SendBytes(json, 0)
}

// 从中心服务器接受数据
func (z *ZMQClient) Recv() ([]byte, error) {
	return z.client.RecvBytes(0)
}

func main() {
	client := NewZMQClient()
	// client.Send(map[string]interface{}{
	// 	"cmd": "ReqFlushConfig",
	// 	"data": map[string]interface{}{
	// 		"configName": "ALL",
	// 	},
	// })

	client.Send(map[string]interface{}{
		"cmd": "ReqFlushConfig",
		"data": map[string]interface{}{
			"configName": "login_bonus.json",
		},
	})

	ret, err := client.Recv()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("ret: %s\n", ret)
}
