package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"mmorpg/zinx/utils"
	"mmorpg/zinx/ziface"
)

type DataPacket struct {
}

func NewDataPacket() *DataPacket {
	return &DataPacket{}
}

func (this *DataPacket) GetHeadLen() uint32 {
	//DataLen uint32 4bytes+Id uint32 4bytes
	return 8
}

/*
	注意：
		使用binary包时，必须使用明确的长度确定的类型，
		例如（bool, int8, uint8, int16, float32, complex64, ...）
		无法处理：int...
*/

func (this *DataPacket) Packet(messages ziface.IMessages) ([]byte, error) {
	//这里使用NewBuffer的原因是，binary.Write的第一个参数要求是io.writer类型的
	buf := bytes.NewBuffer([]byte{})

	//将头部长度以小端字节序写入到buf中（这里大小端的选取不强制，只要读写的方式一致就可）
	if err := binary.Write(buf, binary.LittleEndian, messages.GetMsgLen()); err != nil {
		return nil, err
	}

	//将数据Id写入buf中
	if err := binary.Write(buf, binary.LittleEndian, messages.GetMsgId()); err != nil {
		return nil, err
	}

	//将数据内容写入到buf中
	if err := binary.Write(buf, binary.LittleEndian, messages.GetMsgData()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (this *DataPacket) UnPacket(binaryData []byte) (ziface.IMessages, error) {
	// 创建byte类型的缓冲器，初始值是binaryData
	buf := bytes.NewBuffer(binaryData)
	msg := &Messages{}

	// 从buf中将msg.DataLen读出
	if err := binary.Read(buf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// msg.DataLen大于系统的最大报文长度
	if msg.DataLen > utils.GlobalObject.MaxPackages {
		return nil, errors.New("msg.DataLen > utils.GlobalObject.MaxPackages")
	}
	// 从buf中将msg.Id读出
	if err := binary.Read(buf, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	return msg, nil
}
