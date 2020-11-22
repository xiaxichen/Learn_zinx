package znet

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

// 负责测试datapack拆包 封包的单元测试
func TestDataPack_Pack(t *testing.T) {
	/*
		模拟的服务器
	*/

	// 1 创建 socket
	listen, err := net.Listen("tcp4", "127.0.0.1:8000")
	if err != nil {
		t.Log(err)
		panic(err)
	}
	// 2 从 client 读取数据

	go func() {
		for {
			accept, err := listen.Accept()
			if err != nil {
				t.Log(err)
				panic(err)
			}
			go func(conn net.Conn) {
				pack := NewDataPack()
				for {
					// 第一次从 conn 读，把包的head 读取出来
					headData := make([]byte, pack.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						t.Log(err)
						panic(err)
					}
					msgHead, err := pack.UnPack(headData)
					if err != nil {
						t.Log(err)
						panic(err)
					}
					if msgHead.GetMsgLen() > 0 {
						// 第二次从 conn 读,根据头中的data length 再读取data的内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						// 根据data length的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							t.Log(err)
							panic(err)
						}
						t.Log("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
					}
				}
			}(accept)
		}
	}()
	/*
		模拟客户端
	*/

	dial, err := net.Dial("tcp4", "127.0.0.1:8000")
	if err != nil {
		t.Log(err)
		panic(err)
	}
	pack := NewDataPack()

	m1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	bufferString := bytes.NewBufferString("bytes")
	i := bufferString.Bytes()
	fmt.Println(len(i))
	i3, err := pack.Pack(m1)
	if err != nil {
		t.Log(err)
		panic(err)
	}
	if err != nil {
		t.Log(err)
		panic(err)
	}
	i2 := []byte("zinx")
	fmt.Println(len(i2))
	m2 := &Message{
		Id:      2,
		DataLen: 5,
		Data:    []byte{},
	}
	i4, err := pack.Pack(m2)
	if err != nil {
		t.Log(err)
		panic(err)
	}
	sendData := append(i3, i4...)
	_, err = dial.Write(sendData)
	if err != nil {
		t.Log(err)
		panic(err)
	}
	time.Sleep(5 * time.Second)
}
