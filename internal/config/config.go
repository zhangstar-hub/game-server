package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Config struct {
	Env   EnvConf
	Test  TestConf
	Test2 Test2Conf
}

var config *Config
var mu sync.RWMutex

type EnvConf struct {
	App struct {
		Host string
		Port int
	}
	Mysql struct {
		Address  string
		User     string
		Password string
		Database string
	}
	Redis struct {
		Address string
		DB      int
	}
}

type TestConf struct {
	A int
}

type Test2Conf struct {
	A int
}

func LoadOneConfig(path string, target interface{}) {
	fmt.Printf("load config %s\n", path)
	mu.Lock()
	defer mu.Unlock()
	full_path := fmt.Sprintf("configs/%s", path)
	file, err := os.Open(full_path)
	if err != nil {
		panic("Could not open " + full_path)
	}
	if err := json.NewDecoder(file).Decode(target); err != nil {
		panic("Could not decode " + full_path)
	}
}

func LoadAllConfig() {
	config = &Config{}
	LoadOneConfig("env.json", &config.Env)
	LoadOneConfig("test.json", &config.Test)
	LoadOneConfig("test2.json", &config.Test2)
}

// 获取配置接口
func GetC() *Config {
	mu.RLock()
	defer mu.RUnlock()
	return config
}
