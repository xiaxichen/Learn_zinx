package ziface

/*
	消息管理抽象层
*/
type IMsgHandle interface {
	// 调度执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	// 消息添加具体的处理逻辑
	AddRouter(msgId uint32, router IRouter) error
	// 启动Worker工作池
	StartWorkerPool()
}
