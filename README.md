# Learn_zinx 

# 本框架学习过程是从b站UP 刘丹冰Aceld 视频学习 

### 其中部分函数和逻辑有所改动之后会以此框架为基础进行http等协议的封装。 在此再次感谢 刘丹冰Aceld ！ 

## github仓库地址 
## [Zinx](https://github.com/aceld/zinx)  

### demo Server
```go
package main

import (
	"learn_zinx/zinx/demo"
	"learn_zinx/zinx/znet"
)

// ping test 路由

func main() {
	//1.创建一个server句柄，使用Zinx的api
	server := znet.NewServer("tcp4")

	//2.注册连接的hook回调

	server.SetOnStart(demo.OnStartFunc)
	server.SetOnStop(demo.OnStopFunc)

	//3.给当前zinx框架添加一个router
	server.AddRouter(0, &demo.PingRouterMsg{})
	server.AddRouter(1, &demo.HelloRouterMsg{})

	//4.启动server
	server.Server()
}
```

demo Client
```go
package main

import (
	"io"
	"learn_zinx/zinx/logger"
	"learn_zinx/zinx/znet"
	"net"
	"os"
	"time"
)

/*
模拟客户端
*/

func main() {
	//SigRun()
	logger.Log.Info("client Start!")
	time.Sleep(1 * time.Second)
	// 1 直接连接服务器 ，得到一个conn连接
	dial, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		logger.Log.Error(err)
		os.Exit(1)
	}
	pack := znet.NewDataPack()
	go func() {
		for {
			headData := make([]byte, pack.GetHeadLen())
			_, err = io.ReadFull(dial, headData)
			if err != nil {
				logger.Log.Error(err)
				os.Exit(1)
			}
			msg, err := pack.UnPack(headData)
			if err != nil {
				logger.Log.Error(err)
				os.Exit(1)
			}
			if msg.GetMsgLen() > 0 {
				msgData := make([]byte, msg.GetMsgLen())
				_, err = io.ReadFull(dial, msgData)
				if err != nil {
					logger.Log.Error(err)
					os.Exit(1)
				}
				msg.SetMsgData(msgData)
			}
			logger.Log.Infof("server call back data len: %d, data=%s", msg.GetMsgLen(), string(msg.GetMsgData()))
		}
	}()
	for {
		// 2 连接调用 Write 写数据
		msgPackage := znet.NewMsgPackage(0, []byte("Hello Zinx v 0.9.1"))
		packData, err := pack.Pack(msgPackage)
		if err != nil {
			logger.Log.Error(err)
			os.Exit(1)
		}
		_, err = dial.Write(packData)
		if err != nil {
			logger.Log.Error(err)
			os.Exit(1)
		}
		logger.Log.Info("----------------------------")
		msgPackage = znet.NewMsgPackage(1, []byte("Hello Zinx v 0.9.1"))
		packData, err = pack.Pack(msgPackage)
		if err != nil {
			logger.Log.Error(err)
			os.Exit(1)
		}
		_, err = dial.Write(packData)
		if err != nil {
			logger.Log.Error(err)
			os.Exit(1)
		}
		time.Sleep(3 * time.Second)
	}
}
```
