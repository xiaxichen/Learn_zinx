package znet

import (
	"learn_zinx/zinx/ziface"
)

//实现router之前先嵌入这个BaseRouter基类 然后根据需要对这个基类重写
type BaseRouter struct{}

//这里之所BaseRouter都为空
//是因为有的Router不希望有 PreHandle 和 PostHandle
//所以Router全部继承BaseRouter的好处就是 不需要实现 PreHandle PostHandle

func (router *BaseRouter) PreHandle(request ziface.IRequest) {

}

func (router *BaseRouter) Handle(request ziface.IRequest) {

}

func (router *BaseRouter) PostHandle(request ziface.IRequest) {

}
