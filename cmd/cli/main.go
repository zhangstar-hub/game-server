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

func connTest() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	var data map[string]interface{}
	data = LongJsonTestData1()

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
	conn.Write(message) // 发送消息内容
	fmt.Println("Data sent successfully.")
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
			data = LongJsonTestData1()
		case 2:
			data = LongJsonTestData2()
		}

		// 发送 JSON 数据
		data = map[string]interface{}{
			"cmd":  "cmd",
			"data": data,
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}
		msgLength := make([]byte, 4)
		binary.BigEndian.PutUint32(msgLength, uint32(len(jsonData)))

		fmt.Printf("msgLength: %v %d\n", msgLength, uint32(len(jsonData)))
		message := append(msgLength, jsonData...)
		conn.Write(message) // 发送消息内容
		fmt.Println("Data sent successfully.")
	}

}
