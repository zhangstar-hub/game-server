package ctx

import (
	"my_app/internal/models"
	"net"
)

// 存储服务器所有玩家
var Users = map[int]Ctx{}

type Ctx struct {
	Conn net.Conn // tcp 连接
	Cmd  string   // 当前连接的路由
	User *models.User
}
