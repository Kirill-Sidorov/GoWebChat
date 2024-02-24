package messages

import (
	"fmt"
	"strings"
	"sync"
)

var messages strings.Builder
var lock sync.RWMutex

func GetMessages() string {
	lock.RLock()
	result := messages.String()
	lock.RUnlock()
	return result
}

func AddMessage(userName string, message string) {
	lock.Lock()
	messages.WriteString(fmt.Sprintf("%s: %s\n", userName, message))
	lock.Unlock()
}
