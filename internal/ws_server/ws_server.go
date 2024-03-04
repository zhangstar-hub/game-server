package ws_server

import (
	"fmt"
	"my_app/internal/config"
	"my_app/internal/logger"
	"my_app/internal/src"
	"my_app/internal/zmq_client"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WSServer struct {
	CtxMap      *sync.Map             // 上下文集合
	Listener    *http.Server          // socket监听器
	ZClient     *zmq_client.ZMQClient // zmq_client
	Group       *sync.WaitGroup       // 等待所有连接请求完成
	Stop        chan os.Signal        // 服务关闭信号监听
	CloseFlag   bool                  // 关闭标记
	RoomManager *src.RoomManager      // 房间管理器
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWSServer(Listener *http.Server) *WSServer {
	ctxMap := &sync.Map{}
	roomManager := src.NewRoomManager()

	zClient := zmq_client.NewZMQClient(ctxMap, roomManager)
	go zClient.MessageListener()

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	s := &WSServer{
		CtxMap:      ctxMap,
		Listener:    Listener,
		ZClient:     zClient,
		Group:       &sync.WaitGroup{},
		Stop:        stopChan,
		CloseFlag:   false,
		RoomManager: roomManager,
	}
	go s.ListenSignal()
	return s
}

// 关闭服务
func (s *WSServer) Close() {
	s.Listener.Close()
	s.RoomManager.Close()
	s.CloseFlag = true
}

// 处理连接
func (s *WSServer) handleConnections(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Upgrade error")
		return
	}

	s.Group.Add(1)
	defer s.Group.Done()

	sc := NewWSServerConn(conn)
	defer sc.Close()

	token := uuid.NewString()
	ctx := &src.Ctx{
		Conn:           sc,
		LastActiveTime: time.Now(),
		LastSaveTime:   time.Now(),
		Token:          token,
		ZClient:        s.ZClient,
		RoomManager:    s.RoomManager,
	}
	defer ctx.Close()

	s.CtxMap.Store(token, ctx)
	defer s.CtxMap.Delete(token)

	for !sc.CloseFlag {
		sc.RequestWait()
		data, err, de_err := sc.ReadData()
		if err != nil {
			fmt.Printf("Error reading data: %v, %T\n", err, err)
			return
		}
		if de_err != nil {
			fmt.Printf("Error decoding data: %v, %T\n", de_err, de_err)
			continue
		}
		ret := sc.RequestFunction(ctx, data)
		err = sc.SendData(ret)
		if err != nil {
			fmt.Println("Error sending data:", err)
			continue
		}
	}

}

// 启动信号监听
func (s *WSServer) ListenSignal() {
	defer s.Close()

	fmt.Printf("start listening for signals\n")
	<-s.Stop
	s.Listener.Close()
	fmt.Println("Stop receiving new connections...")
	<-s.Stop

	s.CtxMap.Range(func(key, value interface{}) bool {
		v := value.(*src.Ctx)
		v.Close()
		return true
	})
	fmt.Println("Exiting ...")
}

// 启动服务服务
func StartServer() {
	fmt.Println("WebSocket Server Running")

	conf := config.GetC()
	server := &http.Server{Addr: fmt.Sprintf("%s:%d", conf.Env.App.Host, conf.Env.App.Port)}
	wsServer := NewWSServer(server)
	defer wsServer.Close()

	http.HandleFunc("/", wsServer.handleConnections)
	go wsServer.UserActiveListener()
	go wsServer.AutoSave()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Println(err)
	}
	wsServer.Group.Wait()
	fmt.Println("stop server")
}
