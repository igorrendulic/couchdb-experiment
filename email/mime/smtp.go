package mime

import (
	"fmt"
	"sort"
	"sync"

	"github.com/igorrendulic/couchdb-experiment/email/mime/handler"
)

var (
	handlersMu sync.RWMutex
	handlers   = make(map[string]handler.SmtpHandler)
)

// Register makes a parser available by the provided name.
// If Register is called twice with the same name or if parser is nil,
// it panics.
func Register(name string, handler handler.SmtpHandler) {
	handlersMu.Lock()
	defer handlersMu.Unlock()
	if handler == nil {
		panic("sql: Register driver is nil")
	}
	if _, dup := handlers[name]; dup {
		panic("sql: Register called twice for driver " + name)
	}
	handlers[name] = handler
}

func unregisterAllSmtpHandlers() {
	handlersMu.Lock()
	defer handlersMu.Unlock()
	// For tests.
	handlers = make(map[string]handler.SmtpHandler)
}

// Handlers returns a sorted list of the names of the registered smtp handlers.
func SmtpHandlers() []string {
	handlersMu.RLock()
	defer handlersMu.RUnlock()
	var list []string
	for name := range handlers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

// GetSmtpHandlers returns a handler by name.
func GetSmtpHandler(handlerName string) (handler.SmtpHandler, error) {
	handlersMu.RLock()
	defer handlersMu.RUnlock()
	handler, ok := handlers[handlerName]
	if !ok {
		return nil, fmt.Errorf("sql: unknown parser %q (forgotten import?)", handlerName)
	}
	return handler, nil
}
