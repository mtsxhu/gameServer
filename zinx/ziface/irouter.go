package ziface

type IRouter interface {
	//处理conn业务之前的钩子方法
	PreHandler(request IRequest)
	//在处理conn业务的钩子方法
	Handler(request IRequest)
	//处理conn业务之后的钩子方法
	PostHandler(request IRequest)
}
