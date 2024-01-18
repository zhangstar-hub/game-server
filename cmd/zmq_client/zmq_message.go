package main

// 刷新全部配置
func (z *ZMQClient) FlushAllConfig() {
	z.Send(map[string]interface{}{
		"cmd": "ReqFlushConfig",
		"data": map[string]interface{}{
			"configName": "ALL",
		},
	})
}

// 刷新单个配置
func (z *ZMQClient) FlushOneConfig(configName string) {
	z.Send(map[string]interface{}{
		"cmd": "ReqFlushConfig",
		"data": map[string]interface{}{
			"configName": configName,
		},
	})
}
