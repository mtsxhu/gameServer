package ziface

import "time"

type ZTimer struct {
	// 触发次数,-1表示无限触发
	Count int
	// 定时器ID（唯一标识 ）
	Id int
	// 定时周期
	Cycle time.Duration
	// 定时任务
	Task func(IConnection)
}
type ITimer interface {
	// 添加一个定时器
	AddTimer(*ZTimer)
	// 删除一个定时器
	RemoveTimer(*ZTimer)
	// 获得堆顶的定时器
	GetTop()*ZTimer
	// 删除堆顶定时器
	RemoveTop()*ZTimer
	// 心搏函数
	Tick()
	// 修改触发次数
}

