package tcpx

type Message struct {
	//消息ID
	ID uint8
	//消息的长度
	DataLen uint16
	//消息的内容
	Data []byte
}

// 创建一个message
func NewMessage(id uint8, data []byte) *Message {
	return &Message{
		ID:      id,
		DataLen: uint16(len(data)),
		Data:    data,
	}
}
