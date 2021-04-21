package ziface

import "net"

type IConnection interface {
	//启动链接
	Start()
	//关闭链接
	Stop()
	////获取链接的套接字
	//GetTCPConnection() *net.TCPConn
	//获取链接ID
	GetConnID() int
	//获取远端IP和port
	GetRemoteInfo() net.Addr
	//发送数据
	SendMsg(msgId uint32, data []byte) error
	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string)(interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
	// 添加定时任务
	AddTimer(timer *ZTimer)
}
type HandlerFunc func(*net.TCPConn, []byte, int) error
