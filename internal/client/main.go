package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

func LongJsonTestData1() map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = 1234567890
	data["name"] = "<NAME>"
	data["age"] = 25
	data["address"] = "123 Main St."
	data["city"] = "San Francisco"
	data["state"] = "CA"
	data["zip"] = "94107"
	data["phone"] = "123-456-7890"
	data["email"] = "<EMAIL>"
	data["url"] = "http://www.example.com"
	data["ip"] = "127.0.0.1"
	return data
}

func LongJsonTestData2() map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = 1234567890
	data["name"] = "<NAME>"
	data["age"] = 25
	data["address"] = "123 Main St."
	data["city"] = "San Francisco"
	data["state"] = "CA"
	data["zip"] = "94107"
	data["phone"] = "123-456-7890"
	data["email"] = "<EMAIL>"
	data["url"] = "http://www.example.com"
	data["ip"] = "127.0.0.1"
	for i := 0; i < 1000000; i++ {
		data[fmt.Sprintf("key%d", i)] = i
	}
	return data
}

// 读取消息内容
func readData(conn *net.Conn) ([]byte, error) {
	lenBuffer := make([]byte, 4)
	_, err := (*conn).Read(lenBuffer)
	if err != nil {
		return nil, err
	}
	messageLength := binary.BigEndian.Uint32(lenBuffer)
	var message []byte
	var cap_unm uint32
	for t := messageLength; t > 0; {
		if t > 4096 {
			cap_unm = 4096
		} else {
			cap_unm = t
		}
		new_buffer := make([]byte, cap_unm)
		n, err := (*conn).Read(new_buffer)
		if err != nil {
			return nil, err
		}
		message = append(message, new_buffer[:n]...)
		t -= uint32(n)
	}
	return message, nil
}

func SendTest() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	var n int = 1
	var data map[string]interface{}
	switch n {
	case 1:
		data = map[string]interface{}{
			"cmd": "login",
			"data": map[string]interface{}{
				"name":     "admin",
				"password": "admin",
			},
		}
	case 2:
		data = LongJsonTestData2()
	}

	// 发送 JSON 数据
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
	msgLength := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLength, uint32(len(jsonData)))

	message := append(msgLength, jsonData...)
	_, err = conn.Write(message) // 发送消息内容
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	readData(&conn)
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
				"cmd": "login",
				"data": map[string]interface{}{
					"name":     "admin",
					"password": "admin",
				},
			}
		case 2:
			data = LongJsonTestData2()
		}

		// 发送 JSON 数据
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}
		msgLength := make([]byte, 4)
		binary.BigEndian.PutUint32(msgLength, uint32(len(jsonData)))

		fmt.Printf("msgLength: %v %d\n", msgLength, uint32(len(jsonData)))
		message := append(msgLength, jsonData...)
		_, err = conn.Write(message) // 发送消息内容
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
		ret, _ := readData(&conn)
		fmt.Printf("ret: %s\n", ret)
	}

}
