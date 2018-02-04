package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/log"
)

type Processor struct {
	msgInfo map[int]*MsgInfo
}

type MsgInfo struct {
	msgType       reflect.Type
	msgRouter     *chanrpc.Server
	msgHandler    MsgHandler
	msgRawHandler MsgHandler
}

type MsgHandler func([]interface{})

type MsgRaw struct {
	msgID      int
	msgRawData json.RawMessage
}

func NewProcessor() *Processor {
	p := new(Processor)
	p.msgInfo = make(map[int]*MsgInfo)
	return p
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) Register(msgID int /*msg interface{}*/) int /*string*/ {
	//	msgType := reflect.TypeOf(msg)
	//	if msgType == nil || msgType.Kind() != reflect.Ptr {
	//		log.Fatal("json message pointer required")
	//	}
	//	msgID := msgType.Elem().Name()
	//	if msgID == "" {
	//		log.Fatal("unnamed json message")
	//	}
	if _, ok := p.msgInfo[msgID]; ok {
		log.Fatal("message %v is already registered", msgID)
	}

	i := new(MsgInfo)
	//	i.msgType = msgType
	p.msgInfo[msgID] = i
	return msgID
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRouter(msgID int /*msg interface{}*/, msgRouter *chanrpc.Server) {
	//	msgType := reflect.TypeOf(msg)
	//	if msgType == nil || msgType.Kind() != reflect.Ptr {
	//		log.Fatal("json message pointer required")
	//	}
	//	msgID := msgType.Elem().Name()
	i, ok := p.msgInfo[msgID]
	if !ok {
		for k, v := range p.msgInfo {
			log.Error("k:%v v:%v", k, v)
		}
		log.Fatal("message %v not registered", msgID)
	}

	i.msgRouter = msgRouter
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetHandler(msgID int /*msg interface{}*/, msgHandler MsgHandler) {
	//	msgType := reflect.TypeOf(msg)
	//	if msgType == nil || msgType.Kind() != reflect.Ptr {
	//		log.Fatal("json message pointer required")
	//	}
	//	msgID := msgType.Elem().Name()
	i, ok := p.msgInfo[msgID]
	if !ok {
		log.Fatal("message %v not registered", msgID)
	}

	i.msgHandler = msgHandler
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRawHandler(msgID int /*string*/, msgRawHandler MsgHandler) {
	i, ok := p.msgInfo[msgID]
	if !ok {
		log.Fatal("message %v not registered", msgID)
	}

	i.msgRawHandler = msgRawHandler
}

// goroutine safe
func (p *Processor) Route(msgID int /*msg interface{}*/, userData interface{}, bufMsg []byte) error {
	// raw
	//	if msgRaw, ok := msg.(MsgRaw); ok {
	i, ok := p.msgInfo[msgID /*msgRaw.msgID*/]
	if !ok {
		return fmt.Errorf("message %v not registered", msgID /*msgRaw.msgID*/)
	}
	if i.msgRawHandler != nil {
		//		i.msgRawHandler([]interface{}{msgID /*msgRaw.msgID*/, msgRaw.msgRawData, userData})
	}
	//	}

	// json
	//	msgType := reflect.TypeOf(msg)
	//	if msgType == nil || msgType.Kind() != reflect.Ptr {
	//		return errors.New("json message pointer required")
	//	}
	//	msgID := msgType.Elem().Name()

	//	i, ok := p.msgInfo[msgID]
	//	if !ok {
	//		return fmt.Errorf("message %v not registered", msgID)
	//	}
	//	if i.msgHandler != nil {
	//		i.msgHandler([]interface{}{msgID /*msg*/, userData})
	//	}
	if i.msgRouter != nil {
		i.msgRouter.Go(msgID /*msg*/, userData, bufMsg)
	}
	return nil
}

// goroutine safe
func (p *Processor) Unmarshal(data []byte) (interface{}, error) {
	var m map[int]json.RawMessage /*map[string]json.RawMessage*/
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	if len(m) != 1 {
		return nil, errors.New("invalid json data")
	}

	for msgID, data := range m {
		i, ok := p.msgInfo[msgID]
		if !ok {
			return nil, fmt.Errorf("message %v not registered", msgID)
		}

		// msg
		if i.msgRawHandler != nil {
			return MsgRaw{msgID, data}, nil
		} else {
			msg := reflect.New(i.msgType.Elem()).Interface()
			return msg, json.Unmarshal(data, msg)
		}
	}

	panic("bug")
}

// goroutine safe
func (p *Processor) Marshal(msgID int /*msg interface{}*/) ([][]byte, error) {
	//	msgType := reflect.TypeOf(msg)
	//	if msgType == nil || msgType.Kind() != reflect.Ptr {
	//		return nil, errors.New("json message pointer required")
	//	}
	//	msgID := msgType.Elem().Name()
	if _, ok := p.msgInfo[msgID]; !ok {
		return nil, fmt.Errorf("message %v not registered", msgID)
	}

	// data
	m := map[int]interface{}{msgID: msgID}
	data, err := json.Marshal(m)
	return [][]byte{data}, err
}
