package apis

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"mmorpg/dao"
	"mmorpg/pb"
	"mmorpg/zinx/ziface"
	"mmorpg/zinx/znet"
)

// 玩家登录
type LoginApi struct {
	znet.BaseRouter
}

func (*LoginApi) Handler(request ziface.IRequest){
	//1. 将客户端传来的proto协议解码
	msg:=&pb.LoginPack{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err != nil {
		fmt.Println("Login Unmarshal error ", err)
		return
	}
	//2. TODO 读缓存
	//3. 从数据库中读取用户信息
	password,err:=dao.DBGetUser(msg.UserName)
	if err != nil {
		fmt.Println("DBGetUser error ", err)
		return
	}
	if password!=msg.Password {
		data:=&pb.ErrorPack{ErrInfo: "用户不存在"}
		msg,err:=proto.Marshal(data)
		if err != nil {
			fmt.Println("ErrInfo marshal error ", err)
			return
		}
		request.GetConnection().SendMsg(2103,msg)
	}

}