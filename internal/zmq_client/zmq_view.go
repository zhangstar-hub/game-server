package zmq_client

import (
	"fmt"
	"my_app/internal/config"
	"my_app/internal/src"
	"my_app/internal/utils"
	"sync"
)

func ReqZTest(zClient *ZMQClient, data utils.Dict) {
	jsonData := make(utils.Dict)
	fmt.Printf("jsonData: %#v\n", jsonData)
}

// 把玩家掉线
func ReqZUserExit(zClient *ZMQClient, data utils.Dict) {
	zClient.CtxMap.Range(func(key, value interface{}) bool {
		v := value.(*src.Ctx)
		if v.User != nil && v.User.ID == uint(data["uid"].(float64)) {
			v.Close()
			return false
		}
		return true
	})
}

// 刷新配置
func ReqZFlushConfig(zClient *ZMQClient, data utils.Dict) {
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

// ReqZNotify 通用
func ReqZNotify(zClient *ZMQClient, cmd string, data utils.Dict) {
	to_uid_list := map[uint]struct{}{}
	for _, v := range data["to_uid_list"].([]interface{}) {
		to_uid_list[uint(v.(float64))] = struct{}{}
	}
	form_uid := uint(data["form_uid"].(float64))
	message := data["message"].(utils.Dict)

	zClient.CtxMap.Range(func(key, value any) bool {
		v := value.(*src.Ctx)
		fmt.Printf("v.User.ID: %v\n", v.User.ID)
		fmt.Printf("to_uid_list: %v\n", to_uid_list)
		if _, ok := to_uid_list[v.User.ID]; ok {
			v.Conn.SendData(utils.Dict{
				"cmd": cmd,
				"data": utils.Dict{
					"form_uid": form_uid,
					"message":  message,
				},
			})
		}
		return true
	})
}

// 玩家进入房间
func ReqZEnterRoom(zClient *ZMQClient, data utils.Dict) {
	ReqZNotify(zClient, "ReqEnterRoom", data)
}

// 玩家准备通知
func ReqZRoomReady(zClient *ZMQClient, data utils.Dict) {
	ReqZNotify(zClient, "ReqRoomReady", data)
}

// 玩家叫分
func ReqZCallScore(zClient *ZMQClient, data utils.Dict) {
	ReqZNotify(zClient, "ReqCallScore", data)
}

// 玩家出牌
func ReqZPlayCards(zClient *ZMQClient, data utils.Dict) {
	ReqZNotify(zClient, "ReqPlayCards", data)
}
