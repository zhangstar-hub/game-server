package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"my_app/internal/middleware"
	"my_app/internal/router"
	"my_app/internal/src"
	"my_app/internal/utils"
	"my_app/pkg/protocol.go"
	"my_app/pkg/throttle"
	"net"
	"time"
)

type ServerConn struct {
	conn        net.Conn
	trottleList []throttle.RequestTrottle
	CloseFlag   bool
}

// 创建服务连接
func NewServerConn(conn net.Conn) *ServerConn {
	return &ServerConn{
		conn: conn,
		trottleList: []throttle.RequestTrottle{
			throttle.NewSlidingWindowThrottle(),
			throttle.NewTokenBucketThrottle(),
		},
		CloseFlag: false,
	}
}

// 连接请求限流等待
func (sc *ServerConn) RequestWait() bool {
	ticker := time.Tick(10 * time.Second)
	for _, v := range sc.trottleList {
		select {
		case <-ticker:
			return true
		default:
			if !v.CanRequest() {
				time.Sleep(10 * time.Microsecond)
			}
		}
	}
	return true
}

// 关闭连接资源
func (sc *ServerConn) Close() error {
	sc.conn.Close()
	for _, v := range sc.trottleList {
		v.Close()
	}
	sc.CloseFlag = true
	return nil
}

// 读取消息内容
func (sc *ServerConn) ReadData() (utils.Dict, error, error) {
	lenBuffer := make([]byte, 4)
	_, err := sc.conn.Read(lenBuffer)
	if err != nil {
		return nil, err, nil
	}
	messageLength := binary.BigEndian.Uint32(lenBuffer)
	fmt.Println("messageLength: ", messageLength)

	var message []byte
	var cap_unm uint32
	sc.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	for t := messageLength; t > 0; {
		if t > 4096 {
			cap_unm = 4096
		} else {
			cap_unm = t
		}
		new_buffer := make([]byte, cap_unm)
		n, err := sc.conn.Read(new_buffer)
		if err != nil {
			return nil, err, nil
		}
		message = append(message, new_buffer[:n]...)
		t -= uint32(n)
		sc.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	}
	sc.conn.SetReadDeadline(time.Time{})

	decryptedMessage, err := protocol.Decrypt(message)
	if err != nil {
		return nil, nil, err
	}
	// 解析 JSON 数据
	var data utils.Dict
	if err := json.Unmarshal(decryptedMessage, &data); err != nil {
		return nil, nil, err
	}
	return data, nil, nil
}

// 发送数据
func (sc *ServerConn) SendData(data utils.Dict) (err error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	encryptedMessage, err := protocol.Encrypt(jsonData)
	if err != nil {
		return err
	}

	msgLength := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLength, uint32(len(encryptedMessage)))
	message := append(msgLength, encryptedMessage...)
	sc.conn.Write(message)
	return nil
}

// 执行函数入口
func (sc *ServerConn) RequestFunction(ctx *src.Ctx, data utils.Dict) utils.Dict {
	if _, ok := data["cmd"]; !ok {
		return utils.Dict{
			"error": "invalid command",
		}
	}
	cmd := data["cmd"].(string)
	if _, ok := router.Routers[cmd]; !ok {
		return utils.Dict{
			"error": "invalid command",
		}
	}

	if _, ok := data["data"]; !ok {
		return utils.Dict{
			"error": "invalid data",
		}
	}
	data = data["data"].(utils.Dict)

	ret, err := func() (r utils.Dict, e error) {
		defer func() {
			if err := recover(); err != nil {
				e = err.(error)
				utils.PrintStackTrace()
			}
		}()
		ctx.Cmd = cmd
		for _, f := range middleware.MiddlewareList {
			data = f.BeforeHandle(ctx, data)
		}
		r = router.Routers[cmd](ctx, data)
		for _, f := range middleware.MiddlewareList {
			r = f.AfterHandle(ctx, r)
		}
		return
	}()
	if err != nil {
		return map[string]interface{}{
			"cmd": cmd,
			"data": utils.Dict{
				"error": "server error",
			},
		}
	}
	return utils.Dict{
		"cmd":  cmd,
		"data": ret,
	}
}
