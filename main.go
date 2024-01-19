package main

import (
	"flag"
	"fmt"
	"my_app/internal/config"
	"my_app/internal/db"
	"my_app/internal/server"
	"my_app/internal/utils"
	"os"
)

func main() {
	var progress_id int
	flag.IntVar(&progress_id, "pid", 1, "进程编号")
	flag.Parse()
	defer os.Remove(fmt.Sprintf("logs/app_%d.pid", progress_id))

	config.LoadAllConfig()
	db.InitDB()
	db.InitRedis()
	utils.MirateTable()
	server.StartServer()
}
