package network

type Processor interface {
	// must goroutine safe
	Route(msgID int, userData interface{}, bufMsg []byte) error
	// must goroutine safe
	Unmarshal(data []byte) (interface{}, error)
	// must goroutine safe
	Marshal(msgID int) ([][]byte, error)
}
