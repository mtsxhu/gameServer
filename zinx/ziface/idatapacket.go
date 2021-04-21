package ziface

type IDataPacket interface {
	GetHeadLen() int
	Packet(IMessages) ([]byte, error)
	UnPacket([]byte) (IMessages, error)
}
