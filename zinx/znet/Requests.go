package znet

import (
	"learn_zinx/zinx/ziface"
)

type Requests struct {
	// 已经和客户端建立好的链接
	conn ziface.IConnection
	//客户端请求的数据
	data []byte
}

func (r *Requests) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Requests) GetData() []byte {
	return r.data
}
