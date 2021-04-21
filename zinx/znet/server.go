package znet

import (
	"fmt"
	 "mmorpg/zinx/utils"
	 "mmorpg/zinx/ziface"
	"net"
	"reflect"
)

type Server struct {
	//服务器名称
	Name string
	//服务器IP版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
	//消息管理
	MsgHandler ziface.IMsgHandler
	//连接管理
	ConnManager ziface.IConnManager
	//连接建立之后的hook方法
	OnConnStart func(connection ziface.IConnection)
	//连接销毁之前的hook方法
	OnConnStop func(connection ziface.IConnection)
	// 定时器
	Timer ziface.ITimer
}
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:        name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
		Timer: NewTimePile(),
	}
	fmt.Println("=======>>> NewServer，host:", utils.GlobalObject.Host," port:", utils.GlobalObject.TcpPort)
	return s
}
func (this *Server) Start() {
	fmt.Println("server start ")
	fmt.Printf("Name:%s Port:%d MaxPAckages: %d is start...\n",
		utils.GlobalObject.Name, utils.GlobalObject.TcpPort, utils.GlobalObject.MaxPackages)
	addr, err := net.ResolveTCPAddr(this.IPVersion, fmt.Sprintf("%s:%d", this.IP, this.Port))
	if err != nil {
		fmt.Println("resolve addr err:", err)
	}
	go func() {
		this.MsgHandler.StartWorkerPool()
		listen, err := net.ListenTCP(this.IPVersion, addr)
		if err != nil {
			fmt.Println("listen tcp err:", err)
		}
		connID := 0
		for {
			conn, err := listen.AcceptTCP()
			if this.ConnManager.GetConnTotal() >= utils.GlobalObject.MaxConn {
				//TODO 处理过多的连接
				fmt.Println("too many connection ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			if err != nil {
				fmt.Println("accept tcp err:", err)
			}
			c := NewConnection(this, conn, connID, this.MsgHandler)
			connID++
			go c.Start()
		}
	}()
}
func (this *Server) Stop() {
	this.ConnManager.ClearConn()
}
func (this *Server) Server() {
	this.Start()
	select {}
}
func (this *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	t:=reflect.TypeOf(router)
	fmt.Println("router's type is ",t.Name())
	fmt.Printf("router is %+v\n",router)
	this.MsgHandler.AddRouter(msgId, router)
}

func (this *Server) GetConnManager() ziface.IConnManager {
	return this.ConnManager
}

//添加连接建立之后的hook方法
func (this *Server) SetOnConnStartHook(hookHandler func(connection ziface.IConnection)) {
	this.OnConnStart = hookHandler
	fmt.Println("========>> OnConnStartHook is ready")
}

//添加连接销毁之前的hook方法
func (this *Server) SetOnConnStopHook(hookHandler func(connection ziface.IConnection)) {
	this.OnConnStop = hookHandler
	fmt.Println("========>> OnConnStopHook is ready")

}

//调用连接建立之后的hook方法
func (this *Server) CallOnConnStartHook(connection ziface.IConnection) {
	if this.OnConnStart != nil {
		this.OnConnStart(connection)
	}
}

//调用连接销毁之前的hook方法
func (this *Server) CallOnConnStopHook(connection ziface.IConnection) {
	if this.OnConnStop != nil {
		this.OnConnStop(connection)
	}
}

func (this *Server) AddTimer(timer *ziface.ZTimer){
	this.Timer.AddTimer(timer)
}
func (this *Server) RemoveTimer(timer *ziface.ZTimer){
	this.Timer.RemoveTimer(timer)
}