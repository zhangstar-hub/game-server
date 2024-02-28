package src

import (
	"fmt"
	"my_app/internal/db"
	"my_app/internal/utils"
	"reflect"
	"time"
)

type ZMQInterface interface {
	Send(message map[string]interface{}) (int, error)
	BroastMessage(cmd string, form_uid uint, to_uid_list []uint, message utils.Dict)
	Recv() ([]byte, error)
}

type ConnInterface interface {
	SendData(message map[string]interface{}) error
	Close() error
}

type SaveEntry interface {
	Save() error
}

type Ctx struct {
	Conn           ConnInterface // tcp 连接
	Cmd            string        // 当前处理的命令
	LastActiveTime time.Time     // 上一次存活的时间
	LastSaveTime   time.Time     // 上一次存档的时间
	Token          string        // 登录产生的唯一ID
	ZClient        ZMQInterface  // zmq消息发送器

	User        *User
	LoginBonus  *LoginBonus
	Player      *Player
	RoomManager *RoomManager
}

// 玩家退出清理
func (ctx *Ctx) Close() {
	if ctx.User != nil {
		if ok, r := ctx.RoomManager.GetRoom(ctx.User.RoomID, ctx.User.ID); ok {
			r.LeaveRoom(ctx)
		}
	}
	ctx.SaveAll()
	ctx.Conn.Close()
	if ctx.User != nil {
		SetOffline(ctx.User.ID)
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

// 发送退出消息
func (ctx *Ctx) QuitMessage(uid uint) {
	ctx.ZClient.Send(map[string]interface{}{
		"cmd": "ReqZUserExit",
		"data": map[string]interface{}{
			"uid": uid,
		},
	})
}

// 检测玩家是否在线
func IsOnline(uid uint) bool {
	key := fmt.Sprintf("user:%d", uid)
	ret, _ := db.RedisClient.Exists(key)
	return ret > 0
}

// 标记玩家为在线玩家
func SetOnline(uid uint) {
	key := fmt.Sprintf("user:%d", uid)
	db.RedisClient.Set(key, "1", time.Hour*24*14)
}

// 设置玩家为离线玩家
func SetOffline(uid uint) {
	key := fmt.Sprintf("user:%d", uid)
	db.RedisClient.Delete(key)
}
