package ziface
/*
	路由抽象接口
	路由里的数据都是IRequests

 */

type IRouter interface {
	//在处理conn业务之前的钩子方法 Hook
	PreHandle(request IRequest)
	//在处理conn业务的钩子方法 Hook
	Handle(request IRequest)
	//在处理conn业务的之后钩子方法 Hook
	PostHandle(request IRequest)
}