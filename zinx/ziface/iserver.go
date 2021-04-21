package ziface

type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Server()
	//添加路由
	AddRouter(uint32, IRouter)
	//返回连接管理器
	GetConnManager() IConnManager
	//添加连接建立之后的hook方法
	SetOnConnStartHook(func(connection IConnection))
	//添加连接销毁之前的hook方法
	SetOnConnStopHook(func(connection IConnection))
	//调用连接建立之后的hook方法
	CallOnConnStartHook(connection IConnection)
	//调用连接销毁之前的hook方法
	CallOnConnStopHook(connection IConnection)
 	AddTimer(*ZTimer)
	RemoveTimer(*ZTimer)
}
