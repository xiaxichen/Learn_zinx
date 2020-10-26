package znet

import (
	Log "github.com/sirupsen/logrus"
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

// 从链接读取业务方法
func (c *Connection) StartReader() {
	Log.Infof("Reader Goroutine is running..")
	defer Log.Infof("ConnID = %d Reader is Exit, remote addr is %s", c.ConnID, c.RemoteAddr().String())
	defer c.Stop()
	for {
		//读取client Data 到buffer中,最大512字节
		buf := make([]byte, 512)
		readLength, err := c.Conn.Read(buf)
		if err != nil {
			Log.Errorf("Error ConnID = %d remote addr is %s ,%v", c.ConnID, c.RemoteAddr().String(), err)
			if err.Error() == "EOF" {
				break
			}
			continue
		}

		// 调用当前handle绑定的API
		if err:=c.handleAPI(c.Conn, buf, readLength);err!=nil{
			Log.Errorf("Error handleFunc ConnID = %d remote addr is %s ,%v", c.ConnID, c.RemoteAddr().String(), err)
			break
		}
	}
}

func (c *Connection) Start() {
	Log.Infof("Conn Start().. ConnID = %d", c.ConnID)
	go c.StartReader()
	//todo 启动当前写数据的业务
	//panic("implement me")
}

func (c *Connection) Stop() {
	Log.Infof("Conn Stop().. ConnID = %d", c.ConnID)
	if c.isClose == true {
		return
	}
	c.isClose = true
	c.Conn.Close()
	close(c.ExitChan)
	return
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) bool {
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
