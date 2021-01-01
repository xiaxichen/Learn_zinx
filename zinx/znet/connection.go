package znet

import (
	"errors"
	"fmt"
	"io"
	"learn_zinx/zinx/logger"
	"learn_zinx/zinx/utils"
	"learn_zinx/zinx/ziface"
	"net"
	"sync"
)

/*
连接模块
*/

type Connection struct {
	// 当前Conn隶属于那个server
	TcpServer ziface.IServer

	// 当前socket连接 tcp 套接字
	Conn *net.TCPConn

	// 连接ID
	ConnID uint32

	// 当前的连接状态
	isClose bool

	// 告知当前连接已经退出停止的 channel
	ExitChan chan bool

	// 该连接处理的方法
	MsgHandler ziface.IMsgHandle

	// 连接属性
	property map[string]interface{}

	// 保护连接属性的锁
	propertyLock sync.Mutex

	//无缓冲的管道，用于读写goroutine之间的消息通信
	msgChan chan []byte

	//最大处理字节数
	MaxPackageSize uint32
}

// 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

// 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	if property, ok := c.property[key]; ok {
		return property, nil
	}
	return nil, errors.New(fmt.Sprintf("property no define key:%s", key))
}

// 删除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if _, ok := c.property[key]; ok {
		delete(c.property, key)
	} else {
		logger.Log.Errorf("delete property error no define key:%s", key)
	}

}

// 告知当前连接是否关闭
func (c *Connection) IsClose() bool {
	return c.isClose
}

// 从连接读取业务方法
func (c *Connection) StartReader() {
	logger.Log.Debugf("Reader Goroutine is running..")
	defer logger.Log.Debugf("[ConnID = %d Reader is Exit, remote addr is %s]", c.ConnID, c.RemoteAddr().String())
	defer c.Stop()
	for {
		//读取client Data 到buffer中,最大为配置中的MaxPackageSize
		//buf := make([]byte, c.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	logger.Log.Errorf("Error ConnID = %d remote addr is %s ,%v", c.ConnID, c.RemoteAddr().String(), err)
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
			logger.Log.Errorf("Error ConnID = %d remote addr is %s ,%v", c.ConnID, c.RemoteAddr().String(), err)
			if err.Error() == "EOF" {
				return
			}
		}

		// 拆包,得到msgID 和 msgDatalen 放到消息中
		msg, err := pack.UnPack(headData)
		if err != nil {
			logger.Log.Errorf("UnPack Error ConnID = %d remote addr is %s ,%v", c.ConnID, c.RemoteAddr().String(), err)
			return
		}

		// 根据data length 再次读取Data， 放在msgData中
		if msg.GetMsgLen() > 0 {
			// 第二次从 conn 读,根据头中的data length 再读取data的内容
			data := make([]byte, msg.GetMsgLen())

			// 根据data length的长度再次从io流中读取
			_, err := io.ReadFull(c.Conn, data)
			if err != nil {
				logger.Log.Errorf("Read Conn data for Msg set Error ConnID = %d remote addr is %s ,%v", c.ConnID, c.RemoteAddr().String(), err)
				return
			}
			msg.SetMsgData(data)
		}
		// 从当前 Conn 得到数据 绑定到 Request 中
		req := Requests{
			conn: c,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制,将消息发送给工作池即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 执行注册的路由方法
			// 根据绑定好的MsgId 找到对应处理的handle
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// 开始写入
func (c *Connection) StartWriter() {
	logger.Log.Debug("[Zinx] Write Goroutine is running!")
	defer logger.Log.Infof("[Zinx] %s conn Write is Close!", c.RemoteAddr().String())
	for {

		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				logger.Log.Warnf("Send data error %s", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经退出，此时Writer也要退出
			return
		}
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
		logger.Log.Errorf("Pack error msg Id =%s", msgId)
		return err
	}
	// 发送数据到管道
	c.msgChan <- binaryMsg
	return nil
}
func (c *Connection) Start() {
	logger.Log.Debugf("Conn Start().. ConnID = %d", c.ConnID)
	go c.StartReader()
	// 启动当前写数据的业务
	go c.StartWriter()
	// 按照开发者传递进来，创建连接后需要调用的处理业务
	c.TcpServer.CallOnStart(c)
}

func (c *Connection) Stop() {
	logger.Log.Infof("Conn Stop().. ConnID = %d", c.ConnID)
	// 按照开发者传递进来，停止连接前需要调用的处理业务
	c.TcpServer.CallOnStop(c)
	if c.isClose == true {
		return
	}
	c.isClose = true
	c.Conn.Close()
	// 告知writer 关闭
	c.ExitChan <- true
	// 将当前连接从ConnMgr中摘除
	err := c.TcpServer.GetConnMgr().Remove(c)
	logger.Log.Errorf("Conn Stop error:%s", err)
	close(c.ExitChan)
	close(c.msgChan)
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

// 初始化连接的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, ConnID uint32, MaxPackageSize uint32, handler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:      server,
		Conn:           conn,
		ConnID:         ConnID,
		isClose:        false,
		MsgHandler:     handler,
		ExitChan:       make(chan bool, 1),
		msgChan:        make(chan []byte),
		MaxPackageSize: MaxPackageSize,
		property:       map[string]interface{}{},
	}
	// 将connection 加入到ConnMgr中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}
