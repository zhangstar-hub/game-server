package src

import (
	"fmt"
	"io"
	"my_app/internal/db"
	"my_app/internal/models"
	"reflect"
	"time"
)

type ZMQInterface interface {
	Send(message map[string]interface{}) (int, error)
	Recv() ([]byte, error)
}

type SaveEntry interface {
	Save() error
}

type Ctx struct {
	Conn           io.Closer    // tcp 连接
	Cmd            string       // 当前处理的命令
	LastActiveTime time.Time    // 上一次存活的时间
	LastSaveTime   time.Time    // 上一次存档的时间
	Token          string       // 登录产生的唯一ID
	ZClient        ZMQInterface // zmq消息发送器

	User        *models.UserModel
	LoginBonus  *LoginBonus
	Player      *Player
	RoomManager *RoomManager
}

// 玩家退出清理
func (ctx *Ctx) Close() {
	ctx.Conn.Close()
	ctx.SaveAll()
	if ctx.User != nil {
		ctx.SetOffline(ctx.User.ID)
	}
}

// 保存一个数据
func SaveOne(entity SaveEntry) {
	if reflect.ValueOf(entity).IsNil() {
		return
	}
	if err := entity.Save(); err != nil {
		fmt.Println(err)
	}
}

// 保存全部数据
func (ctx *Ctx) SaveAll() {
	SaveOne(ctx.User)
	SaveOne(ctx.LoginBonus)
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
