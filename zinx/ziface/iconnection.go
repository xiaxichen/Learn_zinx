package ziface

import "net"

type IConnection interface {
	//启动链接 让当前链接开始准备工作
	Start()
	//关闭链接 结束当前链接的工作
	Stop()
	//获取当前链接的绑定 socket conn
	GetTCPConnection() *net.TCPConn
	//获取当前链接模块的链接ID
	GetConnID() uint32
	//获取远程客户端的TCP状态 IP Port
	RemoteAddr() net.Addr
	//发送数据，将数据发送给远程的客户端
	Send(msgId uint32, data []byte) error
}

// 定义一个链接处理业务的方法

type HandleFunc func(*net.TCPConn, []byte, int) error
