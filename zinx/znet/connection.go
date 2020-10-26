package znet

import (
	"learn_zinx/zinx/ziface"
	"net"
)

/*
链接模块
*/

type Connection struct {
	//当前socket链接 tcp 套接字
	Conn *net.TCPConn

	// 链接ID
	ConnID uint32

	//当前的链接状态
	isClose bool

	// 当前链接所绑定的处理方法API
	handleAPI ziface.HandleFunc

	//告知当前链接已经退出停止的 channel
	ExitChan chan bool
}

func (Connection) Start() {
	panic("implement me")
}

func (Connection) Stop() {
	panic("implement me")
}

func (Connection) GetTCPConnection() *net.TCPConn {
	panic("implement me")
}

func (Connection) GetConnID() uint32 {
	panic("implement me")
}

func (Connection) RemoteAddr() *net.Addr {
	panic("implement me")
}

func (Connection) Send(data []byte) bool {
	panic("implement me")
}

// 初始化链接的方法
func NewConnection(conn *net.TCPConn, ConnID uint32, callBack ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    ConnID,
		isClose:   false,
		handleAPI: callBack,
		ExitChan:  make(chan bool, 1),
	}
	return c
}
