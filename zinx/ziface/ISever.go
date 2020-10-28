package ziface


//定义接口方法
type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Server()
	//路由功能：给当前服务器注册一个路由方法，供client的链接处理使用
	AddRouter(router IRouter)
}