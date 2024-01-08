package env

import (
	"encoding/json"
	"fmt"
	"os"
)

type App struct {
	Host string
	Port int
}

type Mysql struct {
	Address  string
	User     string
	Password string
	Database string
}

type Redis struct {
	Address string
	DB      int
}

type Env struct {
	App   App
	Mysql Mysql
	Redis Redis
}

var E *Env

func LoadEnv() {
	fmt.Printf("env init\n")
	file, err := os.Open("env/env.json")
	if err != nil {
		panic("Could not open env.json")
	}
	if err := json.NewDecoder(file).Decode(&E); err != nil {
		panic("Could not decode env.json")
	}
	fmt.Printf("env: %v#\n", E)
}
