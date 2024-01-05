package main

import (
	"fmt"

	"github.com/pebbe/zmq4"
)

func main() {
	client1, _ := zmq4.NewSocket(zmq4.DEALER)
	defer client1.Close()

	_ = client1.Connect("tcp://127.0.0.1:5555")

	// 向服务器发送消息
	_, _ = client1.Send("Client1 message", 0)

	// 接收来自服务器的消息
	msg1, _ := client1.Recv(0)

	fmt.Printf("Received from server by client1: %s\n", msg1)
}
