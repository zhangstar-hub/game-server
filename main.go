package main

import (
	"my_app/internal/server"
	_ "my_app/internal/utils"
)

func main() {
	server.StartServer()
}
