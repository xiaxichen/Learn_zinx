package znet

import (
	"errors"
	"fmt"
	Log "github.com/sirupsen/logrus"
	"learn_zinx/zinx/ziface"
	"net"
	"os"
)

//iServer 的接口实现，定义一个Server的服务器moudel
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的ip
	IP string
	//服务器监听的端口
	Port int
}

// 定义当前客户端所绑定的handle api（目前这个handle是写死的 以后应为自定义）
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	Log.Info("[Conn Handler] CallBackToClient ...")
	if _, err := conn.Write(data); err != nil {
		Log.Errorf("[Conn Handler] write back buf err %v", err)
		return errors.New("CallBackToClient Error")
	}
	return nil
}

func (s *Server) Server() {
	s.Start()
	select {}
}

func (s *Server) Start() {
	var CID uint32
	CID = 0
	Log.Info("[Start] Server listener at IP %s ,Port %d, is starting!", s.IP, s.Port)
	go func() {
		// 1 获取tcp的address
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			Log.Errorf("resolve tcp addt error : %v", err)
			os.Exit(0)
		}
		// 2 监听服务器的地址
		listenIP, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			Log.Errorf(" %v", err)
			os.Exit(0)
		}
		Log.Info("strat Zinx server success ! ", s.Name, " listening...")
		// 3 阻塞的等待客户端连接，处理客户端的请求

		for {
			tcpConn, err := listenIP.AcceptTCP()
			if err != nil {
				Log.Errorf("Accept Error %v", err)
				continue
			}
			//将处理新链接的业务方法
			connection := NewConnection(tcpConn, CID, CallBackToClient)
			CID++
			//启动处理
			go connection.Start()
		}

	}()
}

func (s *Server) Stop() {

}

/*
	初始化Server模块的方法
*/

func NewServer(name, IPVersion, IPAddress string, port int) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: IPVersion,
		IP:        IPAddress,
		Port:      port,
	}
	return s
}
