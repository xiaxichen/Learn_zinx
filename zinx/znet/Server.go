package znet

import (
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

func (s *Server) Server() {
	s.Start()
	select {}
}

func (s *Server) Start() {
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
			// 已经与客户端建立链接，do something 做一个 最大为：512字节长度的回显
			go func() {
				for {
					buf := make([]byte, 512)
					clientIP := tcpConn.RemoteAddr()
					read, err2 := tcpConn.Read(buf)
					if err2 != nil {
						Log.Errorf("%s recv buf err error:%v", clientIP, err2)
						if err2.Error() == "EOF" {
							break
						}
						continue
					}
					// 回显
					if _, err2 := tcpConn.Write(buf[:read]); err2 != nil {
						Log.Errorf("write back buffer error! IP:%s error:%v", clientIP, err2)
						continue
					}

				}
			}()
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
