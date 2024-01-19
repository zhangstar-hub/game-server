package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// 配置文件全局遍历
var config *Config

// 配置文件与存储变量的映射
var ConfigMap map[string]interface{}

// 配置文件锁
var mu sync.RWMutex

// 总配置文件
type Config struct {
	Env        EnvConf
	LoginBonus LoginBonusCFG
}

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
	ZMQCenter struct {
		Address string
	}
}

// 从一个文件中加载配置
func LoadOneConfig(path string, target interface{}, wg *sync.WaitGroup) {
	fmt.Printf("load config %s\n", path)
	wg.Add(1)

	go func() {
		defer wg.Done()
		full_path := fmt.Sprintf("configs/%s", path)
		file, err := os.Open(full_path)
		if err != nil {
			panic("Could not open " + full_path)
		}
		defer file.Close()
		mu.Lock()
		defer mu.Unlock()
		if err := json.NewDecoder(file).Decode(target); err != nil {
			panic("Could not decode " + full_path)
		}
	}()

	ConfigMap[path] = target
}

// 加载所有配置文件
func LoadAllConfig() {
	config = &Config{}
	var wg sync.WaitGroup
	LoadOneConfig("env.json", &config.Env, &wg)
	LoadOneConfig("login_bonus.json", &config.LoginBonus, &wg)
	wg.Wait()
}

// 获取配置接口
func GetC() *Config {
	mu.RLock()
	defer mu.RUnlock()
	return config
}

func init() {
	ConfigMap = make(map[string]interface{})
}
