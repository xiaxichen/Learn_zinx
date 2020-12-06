package znet

import (
	"fmt"
	"learn_zinx/zinx/logger"
	"learn_zinx/zinx/utils"
	"learn_zinx/zinx/ziface"
	"net"
	"os"
)

//iServer 的接口实现，定义一个Server的服务器model
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的ip
	IP string
	//服务器监听的端口
	Port int
	// 路由 当前server的消息管理模块，用来绑定MsgId和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
	// 服务器版本
	ServerVersion string
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
			tcpConn, err := listenIP.AcceptTCP()
			if err != nil {
				logger.Log.Errorf("Accept Error %v", err)
				continue
			}
			//将处理新链接的业务方法
			connection := NewConnection(tcpConn, CID, utils.GlobalObject.MaxPackageSize, s.MsgHandler)
			CID++
			//启动处理
			go connection.Start()
		}

	}()
}

func (s *Server) Stop() {

}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	err := s.MsgHandler.AddRouter(msgId, router)
	if err != nil {
		return
	}
	logger.Log.Info("Add Router success!")
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
	}
	return s
}
