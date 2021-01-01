package ziface

/*
	IRequests接口
	实际上是吧client请求的连接信息和数据包装成到了一个Request中
*/

type IRequest interface {
	//得到当前连接
	GetConnection() IConnection

	//得到请求的消息数据
	GetData() []byte

	//获取数据的Id
	GetId() uint32
}
