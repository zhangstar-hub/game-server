package main

import (
	"encoding/json"
	"fmt"
	"os"

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
	return client.conn.Close()
}

// 接受数据
func (wsConn *WebSocketClient) ReadData() (map[string]interface{}, error, error) {
	_, message, err := wsConn.conn.ReadMessage()
	if err != nil {
		return nil, err, nil
	}
	// message, err = protocol.Decrypt(message)
	// if err != nil {
	// 	return nil, nil, err
	// }
	var data map[string]interface{}
	if err := json.Unmarshal(message, &data); err != nil {
		return nil, nil, err
	}
	return data, nil, nil
}

// 发送数据
func (wsConn *WebSocketClient) SendData(data map[string]interface{}) error {
	message, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// message, err = protocol.Encrypt(message)
	// if err != nil {
	// 	return err
	// }
	err = wsConn.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}
	return nil
}

func SendTest() {
	// 创建 WebSocket 客户端
	serverURL := "ws://localhost:8080/"
	client, err := NewWebSocketClient(serverURL)
	if err != nil {
		fmt.Println("Error creating WebSocket client:", err)
		return
	}
	defer client.Close()

	data := map[string]interface{}{
		"cmd": "ReqKeepAlive",
		"data": map[string]interface{}{
			"name":     "admin",
			"password": "admin",
		},
	}
	// 发送 JSON 数据
	err = client.SendData(data)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	_, err, de_err := client.ReadData()
	if err != nil {
		fmt.Printf("Error reading data: %v, %T\n", err, err)
	}
	if de_err != nil {
		fmt.Printf("Error decoding data: %v, %T\n", de_err, de_err)
	}
}

func ReadListeners(client *WebSocketClient) {
	for {
		ret, err, de_err := client.ReadData()
		if err != nil {
			fmt.Printf("Error reading data: %v, %T\n", err, err)
			os.Exit(0)
		}
		if de_err != nil {
			fmt.Printf("Error decoding data: %v, %T\n", de_err, de_err)
			continue
		}
		r, _ := json.Marshal(ret)
		fmt.Printf("ret: %s\n", r)
	}
}

func main() {
	serverURL := "ws://localhost:8080/"

	// 创建 WebSocket 客户端
	client, err := NewWebSocketClient(serverURL)
	if err != nil {
		fmt.Println("Error creating WebSocket client:", err)
		return
	}
	defer client.Close()
	go ReadListeners(client)
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
		case 5:
			data = map[string]interface{}{
				"cmd": "ReqRoomReady",
				"data": map[string]interface{}{
					"is_ready": true,
				},
			}
		case 6:
			data = map[string]interface{}{
				"cmd":  "ReqEnterRoom",
				"data": map[string]interface{}{},
			}
		case 7:
			data = map[string]interface{}{
				"cmd": "ReqLogin",
				"data": map[string]interface{}{
					"name":     "admin2",
					"password": "admin",
				},
			}
		case 8:
			data = map[string]interface{}{
				"cmd": "ReqLogin",
				"data": map[string]interface{}{
					"name":     "admin3",
					"password": "admin",
				},
			}
		}

		// 发送 JSON 数据
		err = client.SendData(data)
		if err != nil {
			fmt.Println("Error sending data:", err)
			return
		}
	}
}
