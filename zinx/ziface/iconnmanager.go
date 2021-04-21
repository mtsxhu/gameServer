package ziface

type IConnManager interface {
	// 添加连接
	AddConn(IConnection)
	// 删除连接
	RemoveConn(connection IConnection)
	// 查看连接
	GetConnByID(ConnID int) (IConnection, error)
	// 清除所有连接
	ClearConn()
	// 获得连接总数
	GetConnTotal() int
}
