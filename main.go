package main

import (
	"flag"
	"fmt"
	"my_app/internal/config"
	"my_app/internal/db"
	"my_app/internal/models"
	"my_app/internal/ws_server"
	"os"
)

func main() {
	var progress_id int
	flag.IntVar(&progress_id, "pid", 1, "进程编号")
	flag.Parse()
	defer os.Remove(fmt.Sprintf("pids/app_%d.pid", progress_id))

	config.LoadAllConfig()
	db.InitDB()
	db.InitRedis()
	models.MirateTable()
	ws_server.StartServer()
	// server.StartServer()
}
