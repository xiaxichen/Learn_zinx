package ziface

/*
	IRequests接口
	实际上是吧client请求的链接信息和数据包装成到了一个Request中
 */

type IRequest interface {
	//得到当前链接
	GetConnection() IConnection

	//得到请求的消息数据
	GetData() []byte
}