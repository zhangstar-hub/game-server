package ws_server

import (
	"encoding/json"
	"fmt"
	"my_app/internal/config"
	"my_app/internal/logger"
	"my_app/internal/middleware"
	"my_app/internal/router"
	"my_app/internal/src"
	"my_app/internal/utils"
	"my_app/internal/zmq_client"
	"my_app/pkg/throttle"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSServerConn struct {
	conn        *websocket.Conn           // ws 连接
	trottleList []throttle.RequestTrottle // 限流器
}

func NewWSServerConn(conn *websocket.Conn) *WSServerConn {
	return &WSServerConn{
		conn: conn,
		trottleList: []throttle.RequestTrottle{
			throttle.NewTokenBucketThrottle(),
			throttle.NewSlidingWindowThrottle(),
		},
	}
}

// 接受数据
func (wsConn *WSServerConn) readData() (map[string]interface{}, error, error) {
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
func (wsConn *WSServerConn) sendData(data map[string]interface{}) error {
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
				time.Sleep(1 * time.Second)
			}
		}
	}
	return true
}

// 处理连接
func handleConnections(zClient *zmq_client.ZMQClient, group *sync.WaitGroup) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("Upgrade error")
			return
		}
		defer conn.Close()
		group.Add(1)
		defer group.Done()

		token := uuid.NewString()
		ctx := &src.Ctx{
			Conn:           conn,
			LastActiveTime: time.Now(),
			LastSaveTime:   time.Now(),
			Token:          token,
			ZClient:        zClient,
		}
		defer ctx.Close()
		src.Users.Store(token, ctx)

		wsConn := NewWSServerConn(conn)

		for {
			wsConn.RequestWait()
			data, err, de_err := wsConn.readData()
			fmt.Printf("data: %v\n", data)
			if err != nil {
				fmt.Printf("Error reading data: %v, %T\n", err, err)
				return
			}
			if de_err != nil {
				fmt.Printf("Error decoding data: %v, %T\n", de_err, de_err)
				continue
			}
			ret := RequestFunction(ctx, data)
			err = wsConn.sendData(ret)
			if err != nil {
				fmt.Println("Error sending data:", err)
				continue
			}
		}
	}

}

// 执行函数入口
func RequestFunction(ctx *src.Ctx, data utils.Dict) utils.Dict {
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
				fmt.Println("Error:", err)
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
			"error": "server error",
		}
	}
	return map[string]interface{}{
		"cmd":  cmd,
		"data": ret,
	}
}

// 启动信号监听
func ListenSignal(c <-chan os.Signal, server *http.Server) {
	fmt.Printf("start listening for signals\n")
	<-c
	server.Close()
	fmt.Println("Stop receiving new connections...")
	<-c

	src.Users.Range(func(key, value interface{}) bool {
		v := value.(*src.Ctx)
		v.Close()
		return true
	})
	fmt.Println("Exiting ...")
}

// 启动服务服务
func StartServer() {
	fmt.Println("WebSocket Server Running")

	var group sync.WaitGroup
	zClient := zmq_client.NewZMQClient()
	go zClient.MessageListener()

	conf := config.GetC()
	http.HandleFunc("/", handleConnections(zClient, &group))
	server := &http.Server{Addr: fmt.Sprintf("%s:%d", conf.Env.App.Host, conf.Env.App.Port)}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go ListenSignal(c, server)
	go UserActiveListener()
	go AutoSave()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Println(err)
	}
	group.Wait()
	fmt.Println("stop server")
}
