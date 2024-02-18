package server

import (
	"errors"
	"fmt"
	"my_app/internal/config"
	"my_app/internal/src"
	"my_app/internal/zmq_client"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type Server struct {
	CtxMap       *sync.Map             // 上下文集合
	Listener     net.Listener          // socket监听器
	ZClient      *zmq_client.ZMQClient // zmq_client
	Group        *sync.WaitGroup       // 等待所有连接请求完成
	ListenNewReq bool                  // 是否关闭监听
	Stop         chan os.Signal        // 服务关闭信号监听
	CloseFlag    bool                  // 服务关闭标志
}

func NewServer() *Server {
	conf := config.GetC()
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", conf.Env.App.Host, conf.Env.App.Port))
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error resolving address: %s", err)))
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error listening: %s", err)))
	}

	ctxMap := &sync.Map{}
	roomManager := src.NewRoomManager()

	zClient := zmq_client.NewZMQClient(ctxMap, roomManager)
	go zClient.MessageListener()

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	s := &Server{
		CtxMap:       ctxMap,
		Listener:     listener,
		ZClient:      zClient,
		Group:        &sync.WaitGroup{},
		ListenNewReq: true,
		Stop:         stopChan,
		CloseFlag:    false,
	}
	go s.ListenSignal()
	return s
}

// 服务退出清理
func (s *Server) Close() {
	s.Listener.Close()
	s.CloseFlag = true
}

// 处理请求
func (s *Server) handleConnection(conn net.Conn) {
	defer s.Group.Done()
	fmt.Println("Client connected:", conn.RemoteAddr())

	sc := NewServerConn(conn)
	defer sc.Close()

	token := uuid.NewString()
	ctx := &src.Ctx{
		Conn:           sc,
		LastActiveTime: time.Now(),
		LastSaveTime:   time.Now(),
		Token:          token,
		ZClient:        s.ZClient,
	}
	defer ctx.Close()

	s.CtxMap.Store(token, ctx)
	defer s.CtxMap.Delete(token)

	for sc.CloseFlag == false {
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

// 监听连接请求
func (s *Server) HandleServer() {
	for s.ListenNewReq {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			break
		}
		s.Group.Add(1)
		go s.handleConnection(conn)
	}
	s.Group.Wait()
}

// 启动信号监听
func (s *Server) ListenSignal() {
	defer s.Close()
	fmt.Printf("start listening for signals\n")
	<-s.Stop
	s.Listener.Close()
	s.ListenNewReq = false
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
	server := NewServer()
	defer server.Close()
	go server.UserActiveListener()
	go server.AutoSave()
	server.HandleServer()
	fmt.Println("stop server")
}
