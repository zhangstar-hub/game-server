package src

import (
	"fmt"
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

	User *models.User
}

// 玩家退出清理
func (ctx *Ctx) Close() {
	ctx.Conn.Close()
	if ctx.User != nil {
		ctx.SaveAll()
		Users.Delete(ctx.User.ID)
	}
}

// 退出数据保存
func (ctx *Ctx) SaveAll() {
	err := ctx.User.Save()
	if err != nil {
		fmt.Println(err)
	}
}
