package znet

type Messages struct {
	Id      uint32
	DataLen uint32
	Data    []byte
}

func (this *Messages) GetMsgId() uint32 {
	return this.Id
}
func (this *Messages) GetMsgLen() uint32 {
	return this.DataLen
}
func (this *Messages) GetMsgData() []byte {
	return this.Data
}
func (this *Messages) SetMsgId(id uint32) {
	this.Id = id
}
func (this *Messages) SetMsgLen(len uint32) {
	this.DataLen = len
}
func (this *Messages) SetMsgData(data []byte) {
	this.Data = data
}
