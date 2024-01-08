package main

import (
	"my_app/env"
	"my_app/internal/configmanger"
	"my_app/internal/db"
	"my_app/internal/server"
	"my_app/internal/utils"
	"my_app/internal/zmq_client"
)

func main() {
	env.LoadEnv()
	configmanger.LoadCofnig()
	db.InitDB()
	db.InitRedis()
	utils.MirateTable()
	zmq_client.InitZMQClient()
	server.StartServer()
}
