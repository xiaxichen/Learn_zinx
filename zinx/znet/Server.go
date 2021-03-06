package znet

import (
	"fmt"
	"learn_zinx/zinx/logger"
	"learn_zinx/zinx/utils"
	"learn_zinx/zinx/ziface"
	"net"
	"os"
)

// iServer 的接口实现，定义一个Server的服务器model
type Server struct {
	// 服务器的名称
	Name string
	// 服务器绑定的ip版本
	IPVersion string
	// 服务器监听的ip
	IP string
	// 服务器监听的端口
	Port int
	// 路由 当前server的消息管理模块，用来绑定MsgId和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
	// 服务器版本
	ServerVersion string
	// 该server的连接器
	ConnMgr ziface.IConnManager
	// 连接创建前的方法
	OnConnStart func(connection ziface.IConnection)
	// 连接销毁前的方法
	OnConnStop func(connection ziface.IConnection)
}

func (s *Server) Server() {
	s.Start()
	select {}
}

func (s *Server) Start() {
	var CID uint32
	CID = 0
	logger.Log.Infof("[Zinx] Config %+v", utils.GlobalObject)
	logger.Log.Infof("[Zinx] final Config %+v", s)
	logger.Log.Infof("[Start] Server listener at IP %s ,Port %d, is starting!", s.IP, s.Port)
	go func() {
		// 开启消息队列以及woker工作池
		s.MsgHandler.StartWorkerPool()
		// 1 获取tcp的address
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			logger.Log.Errorf("resolve tcp addt error : %v", err)
			os.Exit(0)
		}
		// 2 监听服务器的地址
		listenIP, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			logger.Log.Errorf(" %v", err)
			os.Exit(0)
		}
		logger.Log.Info("strat Zinx server success ! ", s.Name, " listening...")
		// 3 阻塞的等待客户端连接，处理客户端的请求

		for {
			tcpConn, errTcp := listenIP.AcceptTCP()
			if errTcp != nil {
				logger.Log.Errorf("Accept Error %v", err)
				continue
			}
			// 设置最大连接数个数的判断如果超出最大连接数量则关闭此新连接
			if s.ConnMgr.Len() > utils.GlobalObject.MaxConn-1 {
				//给客户端发送一个超出最大连接的连接包
				pack := NewDataPack()
				binaryMsg, err1 := pack.Pack(NewMsgPackage(0, []byte("The maximum number of connections exceeded!")))
				if err1 != nil {
					logger.Log.Errorf("PackError func error! err:%s", err1)
				}
				logger.Log.Warnf("[Server] ConnMgr length greater MaxConn length %d ！！！！ New TcpConn well be close!", utils.GlobalObject.MaxConn)
				_, errWrite := tcpConn.Write(binaryMsg)
				if errWrite != nil {
					logger.Log.Errorf("TcpConn Write Error! err:%s", errWrite)
				}

				errClose := tcpConn.Close()
				if errClose != nil {
					logger.Log.Errorf("TcpConn Close Error! err:%s", errClose)
				}
				continue
			}
			// 将处理新连接的业务方法
			connection := NewConnection(s, tcpConn, CID, utils.GlobalObject.MaxPackageSize, s.MsgHandler)
			CID++
			// 启动处理
			go connection.Start()
		}

	}()
}

func (s *Server) Stop() {
	// 将一些服务器的资源、状态或者一些已经开辟的的连接信息进行停止或者回收
	logger.Log.Infof("[Stop] Zinx server name %s", s.Name)
	s.ConnMgr.Clear()
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	// 注册路由
	err := s.MsgHandler.AddRouter(msgId, router)
	if err != nil {
		logger.Log.Errorf("[Zinx] Add Router Error:err %v", err)
		return
	}
	logger.Log.Info("Add Router success!")
}

// 获取连接管理器
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 注册OnConnStart方法

func (s *Server) SetOnStart(f func(connection ziface.IConnection)) {
	s.OnConnStart = f
}

// 调用OnConnStart方法
func (s *Server) CallOnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		logger.Log.Debugf("-->Call OnConnStart")
		s.OnConnStart(connection)
	}
}

// 注册OnConnStop方法
func (s *Server) SetOnStop(f func(connection ziface.IConnection)) {
	s.OnConnStop = f
}

// 调用OnConnStop方法
func (s *Server) CallOnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		logger.Log.Debugf("-->Call OnConnStop")
		s.OnConnStop(connection)
	}
}

/*
	初始化Server模块的方法
*/

func NewServer(IPVersion string) ziface.IServer {
	s := &Server{
		Name:          utils.GlobalObject.Name,
		IPVersion:     IPVersion,
		ServerVersion: utils.GlobalObject.Version,
		IP:            utils.GlobalObject.Host,
		Port:          utils.GlobalObject.TcpPort,
		MsgHandler:    NewMsgHandle(),
		ConnMgr:       MewConnManager(),
	}
	return s
}
