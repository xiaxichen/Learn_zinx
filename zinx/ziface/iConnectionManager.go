package ziface

/*
连接管理模块抽象层
*/
type IConnManager interface {
	//添加链接
	Add(connection IConnection)
	//删除链接
	Remove(connection IConnection) error
	//根据ConnId获取链接
	Get(connID uint32) (IConnection, error)
	//得到当前链接数
	Len() int
	//清除并终止所有链接
	Clear()
}
