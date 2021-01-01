package ziface

/*
连接管理模块抽象层
*/
type IConnManager interface {
	// 添加连接
	Add(connection IConnection)
	// 删除连接
	Remove(connection IConnection) error
	// 根据ConnId获取连接
	Get(connID uint32) (IConnection, error)
	// 得到当前连接数
	Len() int
	// 清除并终止所有连接
	Clear()
}
