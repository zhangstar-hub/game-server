package main

import (
	"my_app/internal/config"
	"my_app/internal/db"
	"my_app/internal/server"
	"my_app/internal/utils"
	"my_app/internal/zmq_client"
)

func main() {
	config.LoadAllConfig()
	db.InitDB()
	db.InitRedis()
	utils.MirateTable()
	zmq_client.InitZMQClient()
	server.StartServer()
}
