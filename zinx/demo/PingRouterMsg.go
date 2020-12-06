package demo

import (
	Log "github.com/sirupsen/logrus"
	"learn_zinx/zinx/logger"
	"learn_zinx/zinx/ziface"
	"learn_zinx/zinx/znet"
)

// ping test 路由
type PingRouterMsg struct {
	znet.BaseRouter
}

func (router *PingRouterMsg) PreHandle(request ziface.IRequest) {
	logger.Log.Info("Call Router PreHandle")
	err := request.GetConnection().Send(request.GetId(), []byte("before ping..."))
	if err != nil {
		Log.Errorf("Call Router PreHandle err:%s", err)
	}
}

func (router *PingRouterMsg) Handle(request ziface.IRequest) {
	logger.Log.Info("Call Router Handle")
	logger.Log.Infof("recv from client: msgID=%d\tData=%s", request.GetId(), string(request.GetData()))
	err := request.GetConnection().Send(request.GetId(), []byte("ping ping ping"))
	if err != nil {
		Log.Errorf("Call Router Handler err:%s", err)
	}
}

func (router *PingRouterMsg) PostHandle(request ziface.IRequest) {
	logger.Log.Info("Call Router PostHandle")
	err := request.GetConnection().Send(request.GetId(), []byte("after ping..."))
	if err != nil {
		logger.Log.Errorf("Call Router PostHandle err:%s", err)
	}
}
