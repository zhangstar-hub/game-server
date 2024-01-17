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

type ZMQInterface interface {
	Send(message map[string]interface{}) (int, error)
	Recv() ([]byte, error)
}

type Ctx struct {
	Conn           net.Conn     // tcp 连接
	Cmd            string       // 当前处理的命令
	LastActiveTime time.Time    // 上一次存活的时间
	LastSaveTime   time.Time    // 上一次存档的时间
	Token          string       // 登录产生的唯一ID
	ZClient        ZMQInterface // zmq消息发送器

	User          *models.User
	LoginBonusCtx *LoginBonusCtx
}

// 玩家退出清理
func (ctx *Ctx) Close() {
	ctx.Conn.Close()
	ctx.SaveAll()
	Users.Delete(ctx.Token)
	if ctx.User != nil {
		ctx.SetOffline(ctx.User.ID)
	}
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
	ret, _ := db.RedisClient.Exists(key)
	return ret > 0
}

// 标记玩家为在线玩家
func (ctx *Ctx) SetOnline(uid uint) {
	key := fmt.Sprintf("user:%d", uid)
	db.RedisClient.Set(key, "1", time.Hour*24*14)
}

// 设置玩家为离线玩家
func (ctx *Ctx) SetOffline(uid uint) {
	key := fmt.Sprintf("user:%d", uid)
	db.RedisClient.Delete(key)
}

// 发送退出消息
func (ctx *Ctx) QuitMessage(uid uint) {
	ctx.ZClient.Send(map[string]interface{}{
		"cmd": "ReqUserExit",
		"data": map[string]interface{}{
			"uid": uid,
		},
	})
}
