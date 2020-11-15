package demo

import (
	Log "github.com/sirupsen/logrus"
	"learn_zinx/zinx/ziface"
	"learn_zinx/zinx/znet"
)

// ping test 路由
type PingRouter struct {
	znet.BaseRouter
}

func (router *PingRouter) PreHandle(request ziface.IRequest) {
	Log.Info("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
	if err != nil {
		Log.Errorf("Call Router PreHandle err:", err)
	}
}

func (router *PingRouter) Handle(request ziface.IRequest) {
	Log.Info("Call Router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping"))
	if err != nil {
		Log.Errorf("Call ping...ping...ping err:", err)
	}
}

func (router *PingRouter) PostHandle(request ziface.IRequest) {
	Log.Info("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping..."))
	if err != nil {
		Log.Errorf("Call Router PostHandle err:", err)
	}
}
