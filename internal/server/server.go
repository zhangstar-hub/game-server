package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"my_app/internal/config"
	"my_app/internal/middleware"
	"my_app/internal/router"
	"my_app/internal/src"
	"my_app/internal/utils"
	"my_app/internal/zmq_client"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
)

var listenNewReq bool = true

// 处理请求
func handleConnection(conn net.Conn, group *sync.WaitGroup, zClient *zmq_client.ZMQClient) {
	defer group.Done()

	fmt.Println("Client connected:", conn.RemoteAddr())

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

	for {
		CanRequest()
		message, err, de_err := readData(conn)
		if err != nil {
			fmt.Printf("Error reading data: %v, %T\n", err, err)
			return
		}

		if de_err != nil {
			fmt.Printf("Error decoding data: %v, %T\n", de_err, de_err)
			continue
		}

		// 解析 JSON 数据
		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			fmt.Println("Error decoding JSON:", err)
			continue
		}

		ret := RequestFunction(ctx, data)
		err = sendData(conn, ret)
		if err != nil {
			fmt.Println("Error sending data:", err)
			continue
		}
	}
}

// 读取消息内容
func readData(conn net.Conn) ([]byte, error, error) {
	lenBuffer := make([]byte, 4)
	_, err := conn.Read(lenBuffer)
	if err != nil {
		return nil, err, nil
	}
	messageLength := binary.BigEndian.Uint32(lenBuffer)
	fmt.Println("messageLength: ", messageLength)

	var message []byte
	var cap_unm uint32
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	for t := messageLength; t > 0; {
		if t > 4096 {
			cap_unm = 4096
		} else {
			cap_unm = t
		}
		new_buffer := make([]byte, cap_unm)
		n, err := conn.Read(new_buffer)
		if err != nil {
			return nil, err, nil
		}
		message = append(message, new_buffer[:n]...)
		t -= uint32(n)
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	}
	conn.SetReadDeadline(time.Time{})

	decryptedMessage, err := decrypt(message)
	if err != nil {
		return nil, nil, err
	}
	return decryptedMessage, nil, nil
}

// 发送数据
func sendData(conn net.Conn, data map[string]interface{}) (err error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	encryptedMessage, err := Encrypt(jsonData)
	if err != nil {
		return err
	}

	msgLength := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLength, uint32(len(encryptedMessage)))
	message := append(msgLength, encryptedMessage...)
	conn.Write(message)

	return nil
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

// 监听连接请求
func HandleServer(tcp *net.TCPListener) {
	var group sync.WaitGroup
	zClient := zmq_client.NewZMQClient()
	go zClient.MessageListener()

	for listenNewReq {
		conn, err := tcp.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			break
		}
		group.Add(1)
		go handleConnection(conn, &group, zClient)
	}
	group.Wait()
}

// 启动信号监听
func ListenSignal(c <-chan os.Signal, listener *net.TCPListener) {
	fmt.Printf("start listening for signals\n")
	<-c
	listener.Close()
	listenNewReq = false
	fmt.Println("Stop receiving new connections...")
	<-c

	src.Users.Range(func(key, value interface{}) bool {
		v := value.(*src.Ctx)
		v.Close()
		return true
	})
	fmt.Println("Exiting ...")
	os.Exit(0)
}

// 启动服务服务
func StartServer() {
	conf := config.GetC()
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", conf.Env.App.Host, conf.Env.App.Port))
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go ListenSignal(c, listener)
	go UserActiveListener()
	go AutoSave()
	HandleServer(listener)
	fmt.Println("stop server")
}
