package main

import (
	"encoding/json"
	"fmt"
	"my_app/pkg/protocol.go"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	conn *websocket.Conn
}

func NewWebSocketClient(url string) (*WebSocketClient, error) {
	// 使用 gorilla/websocket 进行连接
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return &WebSocketClient{conn: conn}, nil
}

func (client *WebSocketClient) Close() error {
	fmt.Printf("close")
	return client.conn.Close()
}

// 接受数据
func (wsConn *WebSocketClient) readData() (map[string]interface{}, error, error) {
	_, buffer, err := wsConn.conn.ReadMessage()
	if err != nil {
		return nil, err, nil
	}
	decryptedMessage, err := protocol.Decrypt(buffer)
	if err != nil {
		return nil, nil, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(decryptedMessage, &data); err != nil {
		return nil, nil, err
	}
	return data, nil, nil
}

// 发送数据
func (wsConn *WebSocketClient) sendData(data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	encryptedMessage, err := protocol.Encrypt(jsonData)
	if err != nil {
		return err
	}
	err = wsConn.conn.WriteMessage(websocket.BinaryMessage, encryptedMessage)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	serverURL := "ws://localhost:8080/ws"

	// 创建 WebSocket 客户端
	client, err := NewWebSocketClient(serverURL)
	if err != nil {
		fmt.Println("Error creating WebSocket client:", err)
		return
	}
	defer client.Close()

	var n int
	for {
		fmt.Println("please input your choice:")
		fmt.Scanln(&n)
		var data map[string]interface{}
		switch n {
		case 1:
			data = map[string]interface{}{
				"cmd": "ReqLogin",
				"data": map[string]interface{}{
					"name":     "admin",
					"password": "admin",
				},
			}
		case 2:
			data = map[string]interface{}{
				"cmd": "ReqTest",
				"data": map[string]interface{}{
					"test": "test",
				},
			}
			for i := 0; i < 10000; i++ {
				data["data"].(map[string]interface{})[string(i)] = i
			}
		case 3:
			d := map[string]interface{}{
				"coin": 1,
			}
			data = map[string]interface{}{
				"cmd":  "ReqAddCoin",
				"data": d,
			}
		case 4:
			data = map[string]interface{}{
				"cmd":  "ReqZmqTest",
				"data": map[string]interface{}{},
			}
		}
		// 发送 JSON 数据
		err = client.sendData(data)
		if err != nil {
			fmt.Println("Error sending data:", err)
			return
		}

		ret, err, de_err := client.readData()
		if err != nil {
			fmt.Printf("Error reading data: %v, %T\n", err, err)
			return
		}
		if de_err != nil {
			fmt.Printf("Error decoding data: %v, %T\n", de_err, de_err)
			continue
		}
		fmt.Printf("ret: %+v\n", ret)
	}
}