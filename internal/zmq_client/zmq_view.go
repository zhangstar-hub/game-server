package zmq_client

import (
	"fmt"
	"my_app/internal/config"
	"my_app/internal/src"
	"my_app/internal/utils"
	"sync"
)

func ReqTest(data utils.Dict) {
	jsonData := make(utils.Dict)
	fmt.Printf("jsonData: %#v\n", jsonData)
}

// 把玩家掉线
func ReqUserExit(data utils.Dict) {
	src.Users.Range(func(key, value interface{}) bool {
		v := value.(*src.Ctx)
		if v.User != nil && v.User.ID == uint(data["uid"].(float64)) {
			v.Close()
			return false
		}
		return true
	})
}

// 刷新配置
func ReqFlushConfig(data utils.Dict) {
	configName := data["configName"].(string)
	if configName == "ALL" {
		config.LoadAllConfig()
	} else {
		if _, ok := config.ConfigMap[configName]; ok {
			var wg sync.WaitGroup
			config.LoadOneConfig(configName, config.ConfigMap[configName], &wg)
			wg.Wait()
		}
	}
}
