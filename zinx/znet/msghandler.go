package znet

import (
	"fmt"
	 "mmorpg/zinx/utils"
	 "mmorpg/zinx/ziface"
	"reflect"
)

type MsgHandler struct {
	Handlers       map[uint32]ziface.IRouter
	WorkerPoolSize int                     //业务工作Worker池的数量
	TaskQueue      []chan ziface.IRequest //Worker负责取任务的消息队列
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Handlers:       make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}

}
func (this *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := this.Handlers[msgId]; ok {
		fmt.Println("handler is exits")
		return
	}
	this.Handlers[msgId] = router
}

func (this *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := this.Handlers[request.GetMsgId()]
	if !ok {
		fmt.Println("该消息没有定义路由")
		return
	}
	t:=reflect.TypeOf(handler)
	fmt.Println("该消息定义的路由为：",t.Name())
	fmt.Printf("router is %+v\n",handler)
	fmt.Println("DoMsgHandler recive data: ",string(request.GetData()))
	handler.PreHandler(request)
	handler.Handler(request)
	handler.PostHandler(request)
}

//启动worker工作池
func (this *MsgHandler) StartWorkerPool() {
	fmt.Println("======> [ WorkerPool ] init ok")
	for i := 0; i < int(this.WorkerPoolSize); i++ {
		//给消息队列分配内存
		this.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动worker
		go this.StartOneWorker(i, this.TaskQueue[i])
	}
}

//启动一个worker工作流程
func (this *MsgHandler) StartOneWorker(workerID int, workerQueue chan ziface.IRequest) {
	fmt.Println("======> [ WorkerPool ] run ok",workerID)
	for true {
		select {
		case request := <-workerQueue:
			fmt.Println("从消息队列中读取数据")
			this.DoMsgHandler(request)
		}
	}
}

//发送消息给消息队列
func (this *MsgHandler) SendMsgToTaskQueue(req ziface.IRequest) {
	workerID := req.GetConnection().GetConnID() % this.WorkerPoolSize
	this.TaskQueue[workerID] <- req
}
