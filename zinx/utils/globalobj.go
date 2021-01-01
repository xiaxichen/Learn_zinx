package utils

import (
	"encoding/json"
	"io/ioutil"
	"learn_zinx/zinx/logger"
	"learn_zinx/zinx/ziface"
	"os"
)

/*
	存储有关zinx框架的全局参数，供其他模块使用
	大部分参数交由用去配置
*/

type GlobalObj struct {
	/*
		server
	*/
	TcpServer ziface.IServer //当前zinx 全局的Server对象
	Host      string         //监听地址
	TcpPort   int            //服务器端口
	Name      string         //服务器名称
	/*
		zinx
	*/
	Version             string //当去Zinx版本号
	MaxConn             int    //当前主机最大连接数
	MaxPackageSize      uint32 //当前zinx框架数据包的最大值
	WorkerPoolSize      uint32 //当前业务工作池数量Goroutine 为0则为关闭工作池模式
	MaxWorkerTaskLength uint32 //zinx 允许用户最多开辟多少个Worker
}

/*
	从用户配置文件加载
*/
func (g GlobalObj) Reload(file string) {
	readFile, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Log.Error("reload config file fail！err:", err)
		logger.Log.Warn("reload default config")
		return
	}
	err = json.Unmarshal(readFile, &GlobalObject)
	if err != nil {
		logger.Log.Error("reload config file fail！err:", err)
		logger.Log.Warn("reload default config")
	}
	logger.Log.Info("reload config success!")
	logger.Log.Info(GlobalObject)
}

/*
	定义一个全局的对外对象GlobalObj
*/

var GlobalObject *GlobalObj

/*
提供一个初始化方法，初始化当前的GlobalObj
*/
func init() {
	//如果配置文件没有加载就是默认配置
	GlobalObject = &GlobalObj{
		TcpServer:           nil,
		Host:                "0.0.0.0",
		TcpPort:             9000,
		Name:                "ZinxServerApp",
		Version:             "V0.4",
		MaxConn:             2000, // 实际连接数 从0开始计数 如最大连接为2000 则实际上最大数值为1999
		MaxPackageSize:      4096,
		WorkerPoolSize:      10, // worker工作池的消息队列个数
		MaxWorkerTaskLength: 1024, // 每个worker对应的消息队列的任务的数量最大值
	}
	// 应该尝试从conf/zinx.json加载一些用户自定义的参数

	nowPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	GlobalObject.Reload(nowPath + "/zinx/conf/Zinx.json")
}
