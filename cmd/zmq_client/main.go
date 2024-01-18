package main

import "fmt"

func main() {
	client := NewZMQClient()
	// client.FlushOneConfig("login_bonus.json")
	client.FlushAllConfig()
	ret, err := client.Recv()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("ret: %s\n", ret)
}
