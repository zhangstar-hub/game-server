package zmq_client

// 退出消息
func QuitMessage(uid uint) {
	ZClient.Send(map[string]interface{}{
		"cmd": "ReqUserExit",
		"data": map[string]interface{}{
			"uid": uid,
		},
	})
}
