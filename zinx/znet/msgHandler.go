package znet

import (
	"errors"
	"learn_zinx/zinx/logger"
	"learn_zinx/zinx/utils"
	"learn_zinx/zinx/ziface"
)

type MsgHandle struct {
	// 存放每个MsgId 所对应的处理方式
	Apis map[uint32]ziface.IRouter
	// 负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作worker池的worker数量
	WorkerPoolSize uint32
	// 消息队列的最大长度
	MaxWorkerTaskLen uint32
}

// 执行注册的路由
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// Requests 中拿到msgID
	handle, ok := mh.Apis[request.GetId()]
	if !ok {
		logger.Log.Errorf("Api msgId not defind  msgId=%d", request.GetId())
		return
	}
	handle.PreHandle(request)
	handle.Handle(request)
	handle.PostHandle(request)
}

// 根据msgId注册路由
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) error {
	// 判断当前id是否被注册了 如果被注册就返回一个异常
	if _, ok := mh.Apis[msgId]; ok {
		// id 已经注册
		logger.Log.Warnf("Resgistered Api ,msgId=%d", msgId)
		return errors.New("resgistered Api !")
	}
	logger.Log.Infof("Add Api MsgId=%d ; success!", msgId)
	mh.Apis[msgId] = router
	return nil
}

// 新建 MsgHandler
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:             make(map[uint32]ziface.IRouter),
		WorkerPoolSize:   utils.GlobalObject.WorkerPoolSize,
		TaskQueue:        make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		MaxWorkerTaskLen: utils.GlobalObject.MaxWorkerTaskLength,
	}
}

// 启动一个Worker工作池 开启工作池的动作，一个zinx框架只能有一个工作池
func (mh *MsgHandle) StartWorkerPool() {
	// 根据WorkerPoolSize 分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 1.当前的 worker对应的channel消息队列 开辟空间 worker id为i
		mh.TaskQueue[i] = make(chan ziface.IRequest, mh.MaxWorkerTaskLen)
		// 2.启动当前的worker，阻塞等待消息从channel传递进来
		go mh.StartWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个工作流程
func (mh *MsgHandle) StartWorker(workerId int, taskQueue chan ziface.IRequest) {
	logger.Log.Debugf("Worker ID=%d is Started!", workerId)
	// 不断的阻塞等待消息对应的队列消息
	for {
		select {
		// 如果有消息进入，出列的就是一个客户端的Request，执行当前request 绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 将消息提交到消息队列，由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1 消息平均分配个不同的worker（根据客户端建立的ConnId来进行分配）
	wokerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	logger.Log.Debugf("Add ConnId=%d ;requests ID=%d ;To WorkerID=%d", request.GetConnection().GetConnID(),
		request.GetId(), wokerId)
	// 2 将消息发送给对应的worker的TaskQueue
	mh.TaskQueue[wokerId] <- request
}
