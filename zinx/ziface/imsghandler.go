package ziface

type IMsgHandler interface {
	AddRouter(uint32, IRouter)
	DoMsgHandler(IRequest)
	StartWorkerPool()
	SendMsgToTaskQueue(IRequest)
}
