package znet

import (
	"errors"
	Log "github.com/sirupsen/logrus"
	"learn_zinx/zinx/ziface"
)

type MsgHandle struct {
	// 存放每个MsgId 所对应的处理方式
	Apis map[uint32]ziface.IRouter
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// Requests 中拿到msgID
	handle, ok := mh.Apis[request.GetId()]
	if !ok {
		Log.Error("Api msgId not defind  msgId=%d", request.GetId())
		return
	}
	handle.PreHandle(request)
	handle.Handle(request)
	handle.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) error {
	// 判断当前id是否被注册了 如果被注册就返回一个异常
	if _, ok := mh.Apis[msgId]; ok {
		// id 已经注册
		Log.Warnf("Resgistered Api ,msgId=%d", msgId)
		return errors.New("resgistered Api !")
	}
	Log.Infof("Add Api MsgId=%d ; success!", msgId)
	mh.Apis[msgId] = router
	return nil
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}
