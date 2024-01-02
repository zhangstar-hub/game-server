package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Client connected:", conn.RemoteAddr())

	for {
		lenBuffer := make([]byte, 4)
		_, err := conn.Read(lenBuffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client closed:", conn.RemoteAddr())
				return
			}
			fmt.Println("Error reading message length:", err)
			return
		}
		messageLength := binary.BigEndian.Uint32(lenBuffer)

		// 读取消息内容
		var message []byte
		for t := messageLength; t > 0; {
			var cap_unm uint32
			if t > 4096 {
				cap_unm = 4096
			} else {
				cap_unm = t
			}
			new_buffer := make([]byte, cap_unm)
			n, err := conn.Read(new_buffer)
			if err != nil {
				fmt.Println("Error reading message:", err)
				return
			}
			message = append(message, new_buffer[:n]...)
			t -= uint32(n)
		}

		// 解析 JSON 数据
		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			fmt.Println("Error decoding JSON:", err)
			continue
		}
		fmt.Println("received JSON data: ", messageLength)
		// fmt.Println("Received JSON data:", data)
	}
}

func HandleServer(tcp *net.TCPListener) {
	for {
		conn, err := tcp.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}
		go handleConnection(conn)
	}
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		os.Exit(1)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	HandleServer(listener)
	fmt.Println("Server started. Waiting for connections...")
}
