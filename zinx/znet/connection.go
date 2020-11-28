package znet

import (
	"errors"
	Log "github.com/sirupsen/logrus"
	"io"
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

	//告知当前链接已经退出停止的 channel
	ExitChan chan bool

	//该链接处理的方法
	Router ziface.IRouter

	//最大处理字节数
	MaxPackageSize uint32
}

// 从链接读取业务方法
func (c *Connection) StartReader() {
	Log.Infof("Reader Goroutine is running..")
	defer Log.Infof("ConnID = %d Reader is Exit, remote addr is %s", c.ConnID, c.RemoteAddr().String())
	defer c.Stop()
	for {
		//读取client Data 到buffer中,最大为配置中的MaxPackageSize
		//buf := make([]byte, c.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	Log.Errorf("Error ConnID = %d remote addr is %s ,%v", c.ConnID, c.RemoteAddr().String(), err)
		//	if err.Error() == "EOF" {
		//		break
		//	}
		//	continue
		//}
		// 创建一个 拆包解包的对象
		pack := NewDataPack()

		// 读取客户端的msg Head 二进制流 8字节
		headData := make([]byte, pack.GetHeadLen())
		//c.GetTCPConnection()
		_, err := io.ReadFull(c.Conn, headData)
		if err != nil {
			Log.Errorf("Error ConnID = %d remote addr is %s ,%v", c.ConnID, c.RemoteAddr().String(), err)
			if err.Error() == "EOF" {
				break
			}
		}

		// 拆包,得到msgID 和 msgDatalen 放到消息中
		msg, err := pack.UnPack(headData)
		if err != nil {
			Log.Errorf("UnPack Error ConnID = %d remote addr is %s ,%v", c.ConnID, c.RemoteAddr().String(), err)
			break
		}

		// 根据datalen 再次读取Data， 放在msgData中
		if msg.GetMsgLen() > 0 {
			// 第二次从 conn 读,根据头中的data length 再读取data的内容
			data := make([]byte, msg.GetMsgLen())

			// 根据data length的长度再次从io流中读取
			_, err := io.ReadFull(c.Conn, data)
			if err != nil {
				Log.Errorf("Read Conn data for Msg set Error ConnID = %d remote addr is %s ,%v", c.ConnID, c.RemoteAddr().String(), err)
				return
			}
			msg.SetMsgData(data)
		}
		// 从当前 Conn 得到数据 绑定到 Request 中
		req := Requests{
			conn: c,
			msg:  msg,
		}
		//执行注册的路由方法
		go func(request ziface.IRequest) {
			//从路由找到注册绑定的Conn对应的router调用
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
			Log.Info("----------------------------")
		}(&req)

	}
}

// 提供一个send Msg的方法
func (c *Connection) Send(msgId uint32, data []byte) error {
	if c.isClose == true {
		return errors.New("Connection is Closed! for send !")
	}
	// 将data进行封包 msgDataLen|MsgId|Data
	pack := NewDataPack()
	binaryMsg, err := pack.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		Log.Errorf("Pack error msg Id =%s", msgId)
		return err
	}
	_, err = c.Conn.Write(binaryMsg)
	if err != nil {
		Log.Errorf("msg send error msg Id =%s", msgId)
		return err
	}
	return nil
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

// 初始化链接的方法
func NewConnection(conn *net.TCPConn, ConnID uint32, MaxPackageSize uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:           conn,
		ConnID:         ConnID,
		isClose:        false,
		Router:         router,
		ExitChan:       make(chan bool, 1),
		MaxPackageSize: MaxPackageSize,
	}
	return c
}
