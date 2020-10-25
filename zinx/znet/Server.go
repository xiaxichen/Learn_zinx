package znet

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
	panic("implement me")
}

func (s *Server) Start() {

}

func (s *Server) Stop() {

}
