package ziface

import "net"

type IConnection interface {
	// 启动连接 让当前连接开始准备工作
	Start()
	// 关闭连接 结束当前连接的工作
	Stop()
	// 获取当前连接的绑定 socket conn
	GetTCPConnection() *net.TCPConn
	// 获取当前连接模块的连接ID
	GetConnID() uint32
	// 获取远程客户端的TCP状态 IP Port
	RemoteAddr() net.Addr
	// 发送数据，将数据发送给远程的客户端
	Send(msgId uint32, data []byte) error
	// 连接是否关闭
	IsClose() bool
	// 开启写
	StartWriter()
	// 开启读
	StartReader()
	// 设置连接属性
	SetProperty(key string, value interface{})
	// 获取连接属性
	GetProperty(key string) (interface{}, error)
	// 移除连接属性
	RemoveProperty(key string)
}

// 定义一个连接处理业务的方法

type HandleFunc func(*net.TCPConn, []byte, int) error
