package gate

import (
	"net"
	//	"reflect"
	"time"

	"fmt"

	//	"io"
	"bytes"
	"encoding/binary"

	"server/msg"

	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/network"
)

type Gate struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processor       network.Processor
	AgentChanRPC    *chanrpc.Server

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LenMsgID     int
	LittleEndian bool
}

func (gate *Gate) Run(closeSig chan bool) {
	var wsServer *network.WSServer
	if gate.WSAddr != "" {
		wsServer = new(network.WSServer)
		wsServer.Addr = gate.WSAddr
		wsServer.MaxConnNum = gate.MaxConnNum
		wsServer.PendingWriteNum = gate.PendingWriteNum
		wsServer.MaxMsgLen = gate.MaxMsgLen
		wsServer.HTTPTimeout = gate.HTTPTimeout
		wsServer.CertFile = gate.CertFile
		wsServer.KeyFile = gate.KeyFile
		wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go(msg.MsgID_NewAgent, a)
				//				gate.AgentChanRPC.Go("NewAgent", a)
			}
			return a
		}
	}

	var tcpServer *network.TCPServer
	if gate.TCPAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingWriteNum = gate.PendingWriteNum
		tcpServer.LenMsgLen = gate.LenMsgLen
		tcpServer.MaxMsgLen = gate.MaxMsgLen
		tcpServer.LittleEndian = gate.LittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go(msg.MsgID_NewAgent, a)
				//				gate.AgentChanRPC.Go("NewAgent", a)
			}
			return a
		}
	}

	if wsServer != nil {
		wsServer.Start()
	}
	if tcpServer != nil {
		tcpServer.Start()
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
	if tcpServer != nil {
		tcpServer.Close()
	}
}

func (gate *Gate) OnDestroy() {}

type agent struct {
	conn     network.Conn
	gate     *Gate
	userData interface{}
}

//func (server *TCPServer) run() {
//	server.wgLn.Add(1)
//	defer server.wgLn.Done()

//	var tempDelay time.Duration
//	for {
//		conn, err := server.ln.Accept()
//		if err != nil {
//			if ne, ok := err.(net.Error); ok && ne.Temporary() {
//				if tempDelay == 0 {
//					tempDelay = 5 * time.Millisecond
//				} else {
//					tempDelay *= 2
//				}
//				if max := 1 * time.Second; tempDelay > max {
//					tempDelay = max
//				}
//				log.Release("accept error: %v; retrying in %v", err, tempDelay)
//				time.Sleep(tempDelay)
//				continue
//			}
//			return
//		}
//		tempDelay = 0

//		server.mutexConns.Lock()
//		if len(server.conns) >= server.MaxConnNum {
//			server.mutexConns.Unlock()
//			conn.Close()
//			log.Debug("too many connections")
//			continue
//		}
//		server.conns[conn] = struct{}{}
//		server.mutexConns.Unlock()

//		server.wgConns.Add(1)

//		tcpConn := newTCPConn(conn, server.PendingWriteNum, server.msgParser)
//		agent := server.NewAgent(tcpConn)
//		go func() {
//			agent.Run()

//			// cleanup
//			tcpConn.Close()
//			server.mutexConns.Lock()
//			delete(server.conns, conn)
//			server.mutexConns.Unlock()
//			agent.OnClose()

//			server.wgConns.Done()
//		}()
//	}
//}

// use by tcp_server.go 如上的注釋代碼
// go func() {
//   agent.Run()
// 消息的分发中心
func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		fmt.Printf("message data: %v", data)
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}
		if a.gate.Processor != nil {
			log.Debug("route message error: %v", a.gate.Processor)

			bufMsgIDLen := data[:a.gate.LenMsgID]
			bufMsg := data[a.gate.LenMsgID:]
			// read len
			bytesBuffer := bytes.NewBuffer(bufMsgIDLen)
			var tmp int32
			if a.gate.LittleEndian {
				binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
			} else {
				binary.Read(bytesBuffer, binary.BigEndian, &tmp)
			}
			msgID := int(tmp)
			fmt.Printf("bufMsgIDLen : %v\n", bufMsgIDLen)
			fmt.Printf("a.gate.LittleEndian : %v\n", a.gate.LittleEndian)
			fmt.Printf("msgID : %v\n", msgID)
			fmt.Printf("bufMsg : %v\n", bufMsg)
			//			// parse len
			//			var msgLen uint32
			//			switch p.lenMsgLen {
			//			case 1:
			//				msgLen = uint32(bufMsgLen[0])
			//			case 2:
			//				if p.littleEndian {
			//					msgLen = uint32(binary.LittleEndian.Uint16(bufMsgLen))
			//				} else {
			//					msgLen = uint32(binary.BigEndian.Uint16(bufMsgLen))
			//				}
			//			case 4:
			//				if p.littleEndian {
			//					msgLen = binary.LittleEndian.Uint32(bufMsgLen)
			//				} else {
			//					msgLen = binary.BigEndian.Uint32(bufMsgLen)
			//				}
			//			}

			//			fmt.Printf("Read msgLen: %v\n", msgLen)
			//			// check len
			//			if msgLen > p.maxMsgLen {
			//				return nil, errors.New("message too long")
			//			} else if msgLen < p.minMsgLen {
			//				return nil, errors.New("message too short")
			//			}

			//			// data
			//			msgData := make([]byte, msgLen)
			//			if _, err := io.ReadFull(conn, msgData); err != nil {
			//				return nil, err
			//			}

			//			// 大端方式 : msgData 四字节的msgID和二进制协议内容
			//			return msgData, nil

			//			err = a.gate.Processor.Route(msg, a)
			//			if err != nil {
			//				log.Debug("route message error: %v", err)
			//				break
			//			}

			//			msg, err := a.gate.Processor.Unmarshal(data)
			//			if err != nil {
			//				log.Debug("unmarshal message error: %v", err)
			//				break
			//			}
			// Processor 為json.go
			err = a.gate.Processor.Route(msgID, a, bufMsg)
			if err != nil {
				log.Debug("route message error: %v", err)
				break
			}
		}
	}
}

func (a *agent) OnClose() {
	if a.gate.AgentChanRPC != nil {
		//		err := a.gate.AgentChanRPC.Call0("CloseAgent", a)
		err := a.gate.AgentChanRPC.Call0(msg.MsgID_CloseAgent, a)
		if err != nil {
			log.Error("chanrpc error: %v", err)
		}
	}
}

func (a *agent) WriteMsg(data ...[]byte /*msgID int*/ /*msg interface{}*/) {
	//	if a.gate.Processor != nil {
	//	data, err := a.gate.Processor.Marshal(msgID)
	//	if err != nil {
	//		log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
	//		return
	//	}
	err := a.conn.WriteMsg(data...)
	if err != nil {
		//		log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
	}
	//	}
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}
