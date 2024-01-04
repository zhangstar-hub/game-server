package main

import (
	_ "my_app/internal/router"
	"my_app/internal/server"
)

func main() {
	server.StartServer()
}
