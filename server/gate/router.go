package gate

import (
	"server/game"
	"server/msg"

	"server/login"
)

///gate模块决定了,某一个消息具体交给内部某一个模块来处理
func init() {

	///指定消息到game模块,模块间使用ChanRPC通信
	msg.Processor.SetRouter(msg.MsgID_Ok, game.ChanRPC)

	msg.Processor.SetRouter(msg.MsgID_SignUp, login.ChanRPC)

	msg.Processor.SetRouter(msg.MsgID_SignIn, login.ChanRPC)

	msg.Processor.SetRouter(msg.MsgID_Register, login.ChanRPC)

	msg.Processor.SetRouter(msg.MsgID_Up, game.ChanRPC)

	msg.Processor.SetRouter(msg.MsgID_Left, game.ChanRPC)

	msg.Processor.SetRouter(msg.MsgID_Right, game.ChanRPC)

	/// 同步给其他用户自己创建角色
	msg.Processor.SetRouter(msg.MsgID_CREATE_ROLE, game.ChanRPC)

	/// 同步用户个人信息在login的个人模块中
	msg.Processor.SetRouter(msg.MsgID_UpLoad, login.ChanRPC)

	/// 联网模式匹配多人消息处理在game模块中
	msg.Processor.SetRouter(msg.MsgID_Match, game.ChanRPC)

	/// 管理员模式,暂时在游戏主逻辑,game.ChanRPC中处理
	msg.Processor.SetRouter(msg.MsgID_Admin, game.ChanRPC)

	/// 用户向其他用户发送窗口消息
	msg.Processor.SetRouter(msg.MsgID_UserMsg, game.ChanRPC)

	/// 用户向服务器发送确认参加匹配模式
	msg.Processor.SetRouter(msg.MsgID_MatchMode, game.ChanRPC)

	msg.Processor.SetRouter(msg.MsgID_Order, game.ChanRPC)

	msg.Processor.SetRouter(msg.MsgID_Finished, game.ChanRPC)
}
