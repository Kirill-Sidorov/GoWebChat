package messages

import (
	"bytes"
	"sync"
)

var newLineSymbol = []byte{'\n'}
var messagesBuffer *bytes.Buffer = bytes.NewBuffer(newLineSymbol)
var lock sync.RWMutex

func GetMessages() []byte {
	lock.RLock()
	result := messagesBuffer.Bytes()
	lock.RUnlock()
	size := len(result) - 1
	return result[:size]
}

func AddMessage(message []byte) {
	lock.Lock()
	messagesBuffer.Write(append(message, newLineSymbol...))
	lock.Unlock()
}
