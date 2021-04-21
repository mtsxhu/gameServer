package main

import (
	"fmt"
	"mmorpg/apis"
	"mmorpg/core"
	"mmorpg/zinx/ziface"
	"mmorpg/zinx/znet"
)

//当客户端建立连接的时候的hook函数
func OnConnecionStart(conn ziface.IConnection)  {
	//创建一个玩家
	player := core.NewPlayer(conn)
	//同步当前的PlayerID给客户端， 走MsgID:1 消息
	player.SyncPid()
	//同步当前玩家的初始化坐标信息给客户端，走MsgID:200消息
	player.BroadCastStartPosition()
	//将当前新上线玩家添加到worldManager中
	core.WorldMgrObj.AddPlayer(player)

	//将该连接绑定属性Pid
	conn.SetProperty("pid", player.Pid)
	// 给当前连接添加一个定时任务
	timer:=&ziface.ZTimer{
		Count: 1,
		Id:    conn.GetConnID(),
		Cycle: 0,
		Task:  nil,
	}
	conn.AddTimer(timer)
	//==============同步周边玩家上线信息，与现实周边玩家信息========
	player.SyncSurrounding()
	//=======================================================

	fmt.Println("=====> Player pidId = ", player.Pid, " arrived ====")
}

//当客户端断开连接的时候的hook函数
func OnConnectionLost(conn ziface.IConnection) {
	//获取当前连接的Pid属性
	pid, _ := conn.GetProperty("pid")
	//根据pid获取对应的玩家对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//触发玩家下线业务
	if pid != nil {
		player.LostConnection()
	}

	fmt.Println("====> Player ", pid , " left =====")

}

func main() {
	//创建服务器句柄
	s := znet.NewServer("test")

	// 注册客户端连接建立和销毁函数
	s.SetOnConnStartHook(OnConnecionStart)
	s.SetOnConnStopHook(OnConnectionLost)
	// 添加路由处理模块
	s.AddRouter(1,&apis.LoginApi{})
	s.AddRouter(2,&apis.WorldChatApi{})		// 聊天
	s.AddRouter(3,&apis.MoveApi{})		   //  移动
	timer:=&ziface.ZTimer{
		Count: 1,
		Id:    0,
		Cycle: 0,
		Task:  nil,
	}
	s.AddTimer(timer)	// 给全部玩家添加的定时任务
	// 启动服务
	s.Server()
}