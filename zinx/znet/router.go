package znet

import (
	 "mmorpg/zinx/ziface"
)

// 默认路由，如果未添加
type BaseRouter struct {
}

//处理conn业务之前的钩子方法
func (this *BaseRouter) PreHandler(request ziface.IRequest) {
}

//在处理conn业务的钩子方法
func (this *BaseRouter) Handler(request ziface.IRequest) {
}

//处理conn业务之后的钩子方法
func (this *BaseRouter) PostHandler(request ziface.IRequest) {
}
