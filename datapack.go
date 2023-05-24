package tcpx

import (
	"encoding/binary"
)

/**
封包，拆包 模块
直接面向TCP连接中的数据流，用于处理TCP粘包
*/

// data开头的2字节是数据的长度,接下来的2字节是数据的id,在接下来是数据的具体内容

// 封包
func Pack(message *Message) []byte {
	res := make([]byte, message.DataLen+3)
	
	//1.将message的id写入res中
	res[0] = message.ID
	
	//2.将datalen写到res中
	binary.BigEndian.PutUint16(res[1:3], message.DataLen)

	//3.将message的内容写到res中
	copy(res[3:], message.Data[:message.DataLen])

	return res
}

// 拆包
func UnPack(binaryData []byte) *Message{
	// read data  id (1 byte) and len (2 bytes)
	
	msg := &Message{
		ID: binaryData[0],
		DataLen: binary.BigEndian.Uint16(binaryData[1:3]),
	}
	msg.Data = make([]byte, msg.DataLen)
	return msg
}