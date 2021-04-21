package znet

import (
	"errors"
	"fmt"
	"io"
	 "mmorpg/zinx/utils"
	 "mmorpg/zinx/ziface"
	"net"
	"sync"
)

type Connection struct {
	// server实例
	TcpServer  ziface.IServer
	// 客户端连接id
	ConnID     int
	// 连接状态
	IsClosed   bool
	// 通知写GO程退出的channel
	ExitChan   chan bool
	// 客户端连接文件描述符
	Conn       *net.TCPConn
	// 消息处理模块
	MsgHandler ziface.IMsgHandler
	// 写go程和对外提供的接口的通信channel
	// 因为框架不知道何时该像connfd写数据，因此需要通过用户自己调用sendMsg像MsgChan中写数据
	// 来启动写GO程
	MsgChan    chan []byte
	// 为什么读GO程不需要ExitChan通知退出，以及MsgChan来阻塞呢
	// 因为connfd是阻塞的，读go程需要从connfd中读取数据，有数据则运行，没有数据则阻塞，因此不需要msgChan。
	// 而读到EOF则代表对端关闭，因此读go程可以直接退出，索引不需要ExitChan来通知退出。

	//链接属性，对用户提供，用于绑定用户数据和当前连接
	property     map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID int, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		ConnID:     connID,
		IsClosed:   false,
		MsgHandler: handler,
		ExitChan:   make(chan bool),
		Conn:       conn,
		MsgChan:    make(chan []byte),
		property:   make(map[string]interface{}), //对链接属性map初始化
	}
	c.TcpServer.GetConnManager().AddConn(c)
	return c
}

// 读go程回调函数
func (this *Connection) StartReader() {
	for {
		//创建封解包对象
		dp := NewDataPacket()

		//获取头部
		headBuf := make([]byte, dp.GetHeadLen())
		fmt.Println(this.Conn.RemoteAddr())
		if _, err := io.ReadFull(this.Conn, headBuf); err != nil {
			// 如果读到EOF，则关闭当前连接
			if err ==io.EOF {
				this.Stop()
			}
			fmt.Println("io.ReadFull(this.Conn, headBuf) err:", err)
			break
		}

		//解包（只解析头部长度（数据长度），不解析数据内容）
		msg, err := dp.UnPacket(headBuf)
		if err != nil {
			fmt.Println("p.UnPacket(headBuf) err:", err)
			break
		}
		fmt.Println("msg data :", string(msg.GetMsgData()))
		// 如果有数据，将数据内容设置到message对象中
		if msg.GetMsgLen() > 0 {
			data := make([]byte, msg.GetMsgLen())
			if _, err = io.ReadFull(this.Conn, data); err != nil {
				fmt.Println("io.ReadFull(this.GetTCPConnection(), data) err:", err)
				break
			}
			msg.SetMsgData(data)
		}
		// 构建请求消息
		req := Request{
			conn: this,
			msg:  msg,
		}
		//如果启动了worker池，则发送消息给消息队列
		if utils.GlobalObject.WorkerPoolSize > 0 {
			fmt.Println("数据被发送给消息队列")
			this.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go this.MsgHandler.DoMsgHandler(&req)
		}
	}
}
// 写go程回调函数
func (this *Connection) StartWriter() {
	fmt.Println("StartWriter...")
	for {
		select {
		case data := <-this.MsgChan:
			if _, err := this.Conn.Write(data); err != nil {
				fmt.Println("send data err:", err)
				return
			}
		case <-this.ExitChan:
			return
		}
	}
}

// 启动链接
func (this *Connection) Start() {
	go this.StartReader()
	go this.StartWriter()
	//调用该连接建立之后的hook方法
	this.TcpServer.CallOnConnStartHook(this)
}

//关闭链接
func (this *Connection) Stop() {
	fmt.Println("关闭连接。。。")
	if this.IsClosed == true {
		return
	}
	this.IsClosed = true
	// 调用该连接销毁前的hook方法
	this.TcpServer.CallOnConnStopHook(this)
	//关闭当前连接文件描述符
	this.Conn.Close()
	this.ExitChan <- true
	// 从连接管理模块中删除该连接
	this.TcpServer.GetConnManager().RemoveConn(this)
	close(this.ExitChan)
	close(this.MsgChan)
}

//获取链接的套接字
//func (this *Connection) GetTCPConnection() *net.TCPConn {
//	return this.Conn
//}

//获取链接ID
func (this *Connection) GetConnID() int {
	return this.ConnID
}

//获取远端IP和port
func (this *Connection) GetRemoteInfo() net.Addr {
	return this.Conn.RemoteAddr()
}

//发送数据
func (this *Connection) SendMsg(msgId uint32, data []byte) error {
	if this.IsClosed == true {
		return errors.New("try to send msg when connection is closed")
	}
	//封包
	dp := NewDataPacket()
	msg := &Messages{
		Id:      msgId,
		DataLen: uint32(len(data)),
		Data:    data,
	}
	binaryMsg, err := dp.Packet(msg)
	if err != nil {
		fmt.Println("dp.Packet(msg) err", err)
	}

	//发送给写模块
	this.MsgChan <- binaryMsg
	return nil
}

// 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok  {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

func (c *Connection)AddTimer(timer *ziface.ZTimer){

}
