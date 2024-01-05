package main

// 这个服务器只用来广播其他服务的消息

import (
	"fmt"

	"github.com/pebbe/zmq4"
)

func main() {
	server, _ := zmq4.NewSocket(zmq4.ROUTER)
	defer server.Close()

	_ = server.Bind("tcp://*:5555")
	// 记录已连接的客户端
	clients := make(map[string]bool)

	for {
		msg, _ := server.RecvMessage(0)

		clientID, message := msg[0], msg[1]
		clients[clientID] = true
		fmt.Printf("Received message from client %s: %s\n", clientID, message)
		if message == "first_message" {
			continue
		}
		for client := range clients {
			_, _ = server.SendMessage(client, message)
		}
	}
}
