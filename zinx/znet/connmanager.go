package znet

import (
	"errors"
	"fmt"
	"mmorpg/zinx/ziface"
	"sync"
)

type ConnManager struct {
	// 存放连接实例和连接id
	connMap sync.Map
}

func NewConnManager() *ConnManager {
	return &ConnManager{}
}

// 添加连接
func (this *ConnManager) AddConn(conn ziface.IConnection) {
	this.connMap.Store(conn.GetConnID(), conn)
	fmt.Println("connection ", conn.GetConnID(), " add to connManager successfuly")
}

// 删除连接
func (this *ConnManager) RemoveConn(conn ziface.IConnection) {
	this.connMap.Delete(conn.GetConnID())
	fmt.Println("connection ", conn.GetConnID(), " delete to connManager successfuly")
}

// 根据连接id获取连接实例
func (this *ConnManager) GetConnByID(ConnID int) (ziface.IConnection, error) {
	conn, ok := this.connMap.Load(ConnID)
	if !ok {
		return nil, errors.New("connection not found")
	}
	return conn.(ziface.IConnection), nil
}

// 清除所有连接
func (this *ConnManager) ClearConn() {
	this.connMap.Range(func(key, value interface{}) bool {
		value.(ziface.IConnection).Stop()
		this.connMap.Delete(value.(ziface.IConnection).GetConnID())
		return true
	})
}

// 获得连接总数
func (this *ConnManager) GetConnTotal() int {
	var sum int
	this.connMap.Range(func(key, value interface{}) bool {
		sum++
		return true
	})
	return sum
}
