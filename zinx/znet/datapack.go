package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"learn_zinx/zinx/logger"
	"learn_zinx/zinx/utils"
	"learn_zinx/zinx/ziface"
)

// 拆包封包的具体模块
type DataPack struct {
}

func (d *DataPack) GetHeadLen() uint32 {
	// Data length uint32 (4字节) + ID uint32(4字节)
	return 8
}

func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	buffer := bytes.NewBuffer([]byte{})

	// 将data length 写入buffer中
	err := binary.Write(buffer, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		goto ERR
	}
	// 将MsgID 写入buffer中
	err = binary.Write(buffer, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		goto ERR
	}
	// 将data数据写入buffer中
	err = binary.Write(buffer, binary.LittleEndian, msg.GetMsgData())
	if err != nil {
		goto ERR
	}

	return buffer.Bytes(), nil
ERR:
	logger.Log.Errorf("Pack error for message:%+v ; err:%s", msg, err)
	return nil, err
}

// 拆包方法 （将包的Head信息读出来） 之后再根据head的信息里的data长度，再进行一次读。
func (d *DataPack) UnPack(bytesData []byte) (ziface.IMessage, error) {
	// 创建一个buffer
	buffer := bytes.NewBuffer(bytesData)
	//初始化msg
	msg := &Message{}

	//读data length
	err := binary.Read(buffer, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		goto ERR
	}
	//读Msg ID
	err = binary.Read(buffer, binary.LittleEndian, &msg.Id)
	if err != nil {
		goto ERR
	}
	//判断data length 是否超出了我们允许的最大长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		err = errors.New("too Large msg data recv!")
		goto ERR
	}

	return msg, err
ERR:
	logger.Log.Errorf("UnPack error for err:%s ; data:%s", err, bytesData)
	return nil, err
}

func NewDataPack() *DataPack {
	return &DataPack{}
}
