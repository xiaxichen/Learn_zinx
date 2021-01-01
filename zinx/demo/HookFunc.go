package demo

import (
	"learn_zinx/zinx/logger"
	"learn_zinx/zinx/ziface"
)

func OnStartFunc(connection ziface.IConnection) {
	logger.Log.Infof("Call OnStartFunc connection id %d", connection.GetConnID())
	connection.Send(200,[]byte("Call OnStartFunc connection BEGIN"))
	connection.SetProperty("hello","world")
}

func OnStopFunc(connection ziface.IConnection) {
	logger.Log.Infof("Call OnStopFunc connection id %d", connection.GetConnID())
	connection.Send(200,[]byte("Call OnStartFunc connection After"))
	logger.Log.Info("hello")
	if property, err := connection.GetProperty("hello");err==nil{
		logger.Log.Info(property)
	}
	connection.RemoveProperty("hello")
	if _, err := connection.GetProperty("hello");err!=nil{
		logger.Log.Warn(err)
	}

}
