package src

import (
	"fmt"
	"my_app/internal/db"
	"my_app/internal/models"
	"net"
	"sync"
	"time"
)

// 存储服务器所有玩家
var Users sync.Map

type Ctx struct {
	Conn           net.Conn  // tcp 连接
	Cmd            string    // 当前处理的命令
	LastActiveTime time.Time // 上一次存活的时间
	LastSaveTime   time.Time // 上一次存档的时间
	Token          string    // 登录产生的唯一ID

	User *models.User
}

// 玩家退出清理
func (ctx *Ctx) Close() {
	ctx.Conn.Close()
	ctx.SaveAll()
	Users.Delete(ctx.Token)
}

// 退出数据保存
func (ctx *Ctx) SaveAll() {
	err := ctx.User.Save()
	if err != nil {
		fmt.Println(err)
	}
}

// 检测玩家是否在线
func (ctx *Ctx) IsOnline(uid uint) bool {
	key := fmt.Sprintf("user:%d", uid)
	client := db.NewRedis()
	ret, _ := client.Exists(key)
	return ret > 0
}
