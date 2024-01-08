package zmq_client

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

var ZClient ZMQClient

func init() {
	client, _ := zmq4.NewSocket(zmq4.DEALER)
	client.SetIdentity(fmt.Sprintf("%s:%d", utils.GetLocalIP(), utils.GetPid()))

	ZClient = ZMQClient{client: client}
	err := client.Connect(ZmqServerAddr)
	if err != nil {
		panic(fmt.Sprintf("无法连接中心服务器：%s", ZmqServerAddr))
	}
	_, err = client.Send("first_message", 0)
	if err != nil {
		panic("向中心服务器发送消息失败")
	}
}

// 向中心服务器发送数据
func (z ZMQClient) Send(data map[string]interface{}) (int, error) {
	json, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	return z.client.Send(string(json), 0)
}

// 从中心服务器接受数据
func (z ZMQClient) Recv() (string, error) {
	return z.client.Recv(0)
}

// 数据接口监听
func MessageListener() {
	for {
		message, err := ZClient.Recv()
		fmt.Printf("MessageListener message: %v\n", message)
		if err != nil {
			fmt.Printf("message recv error %s\n", err)
			continue
		}
		messageMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(message), &messageMap)
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
		ZMQRouters[cmd](data)
	}
}
