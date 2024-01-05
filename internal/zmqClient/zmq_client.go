package zmqClient

import (
	"encoding/json"
	"fmt"

	"github.com/pebbe/zmq4"
)

const (
	ZmqServerAddr = "tcp://127.0.0.1:5555"
)

type zmqClient struct {
	client *zmq4.Socket
}

var ZClient zmqClient

func init() {
	client, _ := zmq4.NewSocket(zmq4.DEALER)
	ZClient = zmqClient{client: client}

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
func (z zmqClient) Send(data map[string]interface{}) (int, error) {
	json, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	return z.client.Send(string(json), 0)
}

// 从中心服务器接受数据
func (z zmqClient) Recv() (string, error) {
	return z.client.Recv(0)
}

// 数据接口监听
func MessageListener() {
	for {
		message, err := ZClient.Recv()
		if err != nil {
			fmt.Printf("数据接受失败\n")
			continue
		}
		fmt.Printf("接收到数据：%s\n", message)
	}
}
