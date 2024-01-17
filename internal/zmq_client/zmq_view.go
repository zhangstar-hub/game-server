package zmq_client

import (
	"fmt"
	"my_app/internal/src"
	"my_app/internal/utils"
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
