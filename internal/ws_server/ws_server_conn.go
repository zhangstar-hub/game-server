package ws_server

import (
	"encoding/json"
	"my_app/internal/middleware"
	"my_app/internal/router"
	"my_app/internal/src"
	"my_app/internal/utils"
	"my_app/pkg/throttle"
	"time"

	"github.com/gorilla/websocket"
)

type WSServerConn struct {
	conn        *websocket.Conn           // ws 连接
	trottleList []throttle.RequestTrottle // 限流器
	CloseFlag   bool                      // 关闭标志
}

func NewWSServerConn(conn *websocket.Conn) *WSServerConn {
	return &WSServerConn{
		conn: conn,
		trottleList: []throttle.RequestTrottle{
			throttle.NewTokenBucketThrottle(),
			throttle.NewSlidingWindowThrottle(),
		},
		CloseFlag: false,
	}
}

// 清理操作
func (c *WSServerConn) Close() error {
	c.conn.Close()
	for _, v := range c.trottleList {
		v.Close()
	}
	c.CloseFlag = true
	return nil
}

// 接受数据
func (wsConn *WSServerConn) ReadData() (map[string]interface{}, error, error) {
	_, message, err := wsConn.conn.ReadMessage()
	if err != nil {
		return nil, err, nil
	}
	// message, err = protocol.Decrypt(message)
	// if err != nil {
	// 	return nil, nil, err
	// }
	var data map[string]interface{}
	if err := json.Unmarshal(message, &data); err != nil {
		return nil, nil, err
	}
	return data, nil, nil
}

// 发送数据
func (wsConn *WSServerConn) SendData(data map[string]interface{}) error {
	message, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// message, err = protocol.Encrypt(message)
	// if err != nil {
	// 	return err
	// }
	err = wsConn.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}
	return nil
}

func (wsConn *WSServerConn) RequestWait() bool {
	ticker := time.Tick(10 * time.Second)
	for _, v := range wsConn.trottleList {
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

// 执行函数入口
func (sc *WSServerConn) RequestFunction(ctx *src.Ctx, data utils.Dict) utils.Dict {
	if _, ok := data["cmd"]; !ok {
		return map[string]interface{}{
			"error": "invalid command",
		}
	}
	cmd := data["cmd"].(string)
	if _, ok := router.Routers[cmd]; !ok {
		return map[string]interface{}{
			"error": "invalid command",
		}
	}

	if _, ok := data["data"]; !ok {
		return map[string]interface{}{
			"error": "invalid data",
		}
	}
	data = data["data"].(map[string]interface{})

	ret, err := func() (r map[string]interface{}, e error) {
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
	return map[string]interface{}{
		"cmd":  cmd,
		"data": ret,
	}
}
