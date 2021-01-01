package znet

import (
	"learn_zinx/zinx/ziface"
)

type Requests struct {
	// 已经和客户端建立好的连接
	conn ziface.IConnection
	//客户端请求的数据
	msg ziface.IMessage
}

func (r *Requests) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Requests) GetData() []byte {
	return r.msg.GetMsgData()
}
func (r *Requests) GetId() uint32 {
	return r.msg.GetMsgId()
}
