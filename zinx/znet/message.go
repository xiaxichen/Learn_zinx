package znet

type Message struct {
	Id      uint32 // 消息ID
	DataLen uint32 // 消息的长度
	Data    []byte // 消息数据

}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

func (m *Message) GetMsgData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(ID uint32) {
	m.Id = ID
}

func (m *Message) SetMsgLen(length uint32) {
	m.DataLen = length
}

func (m *Message) SetMsgData(data []byte) {
	m.Data = data

}
