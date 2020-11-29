package demo

import (
	Log "github.com/sirupsen/logrus"
	"learn_zinx/zinx/ziface"
	"learn_zinx/zinx/znet"
)

// ping test 路由
type HelloRouterMsg struct {
	znet.BaseRouter
}

func (router *HelloRouterMsg) PreHandle(request ziface.IRequest) {
	Log.Info("Call Router PreHandle")
	err := request.GetConnection().Send(request.GetId(), []byte("before hello..."))
	if err != nil {
		Log.Errorf("Call Router PreHandle err:%s", err)
	}
}

func (router *HelloRouterMsg) Handle(request ziface.IRequest) {
	Log.Info("Call Router Handle")
	Log.Infof("recv from client: msgID=%d\tData=%s", request.GetId(), string(request.GetData()))
	err := request.GetConnection().Send(request.GetId(), []byte("hello hello hello"))
	if err != nil {
		Log.Errorf("Call Router Handler err:%s", err)
	}
}

func (router *HelloRouterMsg) PostHandle(request ziface.IRequest) {
	Log.Info("Call Router PostHandle")
	err := request.GetConnection().Send(request.GetId(), []byte("after world..."))
	if err != nil {
		Log.Errorf("Call Router PostHandle err:%s", err)
	}
}
