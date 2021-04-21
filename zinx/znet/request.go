package znet

import "mmorpg/zinx/ziface"

type Request struct {
	// 连接实例
	conn ziface.IConnection
	// 消息实例
	msg  ziface.IMessages
}

// 获取当前请求的连接实例
func (this *Request) GetConnection() ziface.IConnection {
	return this.conn
}
// 获取当前请求消息
func (this *Request) GetData() []byte {
	return this.msg.GetMsgData()
}
// 获取当前请求消息id
func (this *Request) GetMsgId() uint32 {
	return this.msg.GetMsgId()
}
