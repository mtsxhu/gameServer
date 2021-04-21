package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mmorpg/zinx/ziface"
	"os"
)

type GlobalObj struct {
	TCPServer   ziface.IServer
	Host        string
	TcpPort     int
	Name        string
	Version     string
	MaxConn     int
	MaxPackages uint32
	WorkerPoolSize   int //业务工作Worker池的数量
	MaxWorkerTaskLen int //业务工作Worker对应负责，的任务队列中存储的任务最大数量
	Heartbeat	int // 心跳，用于主动与客户端断开连接,精度为秒。-1表示不启动
}

var GlobalObject *GlobalObj

func (this *GlobalObj) Load() {
	data, err := ioutil.ReadFile("./conf/zinx.json")
	if err != nil {
		str, _ := os.Getwd()
		fmt.Println("ReadFile   ", str)
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Host:             "127.0.0.1",
		TcpPort:          9999,
		Name:             "ZinxServerAPP",
		Version:          "V0.4",
		MaxConn:          10000,
		MaxPackages:      4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		Heartbeat: -1,
	}
	GlobalObject.Load()
	fmt.Println("configer is lodding")
}
