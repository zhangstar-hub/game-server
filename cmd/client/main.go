package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"my_app/pkg/protocol.go"
	"net"
	"time"
)

// 读取消息内容
func ReadData(conn net.Conn) ([]byte, error) {
	lenBuffer := make([]byte, 4)
	_, err := conn.Read(lenBuffer)
	if err != nil {
		return nil, err
	}
	messageLength := binary.BigEndian.Uint32(lenBuffer)

	var message []byte
	var cap_unm uint32
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	for t := messageLength; t > 0; {
		if t > 4096 {
			cap_unm = 4096
		} else {
			cap_unm = t
		}
		new_buffer := make([]byte, cap_unm)
		n, err := conn.Read(new_buffer)
		if err != nil {
			return nil, err
		}
		message = append(message, new_buffer[:n]...)
		t -= uint32(n)
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	}
	conn.SetReadDeadline(time.Time{})

	decryptedMessage, err := protocol.Decrypt(message)
	if err != nil {
		return nil, err
	}
	return decryptedMessage, nil
}

// 发送数据
func SendData(conn net.Conn, data map[string]interface{}) (err error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	encryptedMessage, err := protocol.Encrypt(jsonData)
	if err != nil {
		return err
	}

	msgLength := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLength, uint32(len(encryptedMessage)))
	message := append(msgLength, encryptedMessage...)
	_, err = conn.Write(message)
	if err != nil {
		return err
	}
	return nil
}

func SendTest() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	data := map[string]interface{}{
		"cmd": "ReqKeepAlive",
		"data": map[string]interface{}{
			"name":     "admin",
			"password": "admin",
		},
	}
	err = SendData(conn, data)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}
	ReadData(conn)
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

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
		err = SendData(conn, data)
		if err != nil {
			fmt.Println("Error sending data:", err)
			return
		}
		ret, _ := ReadData(conn)
		fmt.Printf("ret: %s\n", ret)
	}

}
