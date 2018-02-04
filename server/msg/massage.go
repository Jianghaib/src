package msg

const (
	MsgID_Ok        = 1
	MsgID_SignUp    = 2
	MsgID_SignIn    = 3
	MsgID_State     = 4
	MsgID_Up        = 5
	MsgID_Right     = 6
	MsgID_Left      = 7
	MsgID_Down      = 8
	MsgID_Command   = 9
	MsgID_UpLoad    = 10
	MsgID_Match     = 11
	MsgID_Admin     = 12
	MsgID_UserMsg   = 13
	MsgID_MatchMode = 14
	MsgID_Order     = 15
	MsgID_Finished  = 16
	// 101 - 300
	MsgID_Register = 101
	MsgID_CREATE_ROLE = 102
	// 301 - 5000
	MsgID_StatusNotValidFile = 603
	// SYSTEM 10001 - 20000
	MsgID_NewAgent   = 10001
	MsgID_CloseAgent = 10002
	MsgID_help       = 10003
	MsgID_cpuprof    = 10004
	MsgID_prof       = 10005
)

//	Processor.Register(&Ok{})
//	Processor.Register(&SignUp{})
//	Processor.Register(&SignIn{})
//	Processor.Register(&State{})
//	Processor.Register(&Up{})
//	Processor.Register(&Right{})
//	Processor.Register(&Left{})
//	Processor.Register(&Down{})
//	Processor.Register(&Command{})
//	Processor.Register(&UpLoad{})
//	Processor.Register(&Match{})
//	Processor.Register(&Admin{})
//	Processor.Register(&UserMsg{})
//	Processor.Register(&MatchMode{})
//	Processor.Register(&Order{})
//	Processor.Register(&Finished{})
