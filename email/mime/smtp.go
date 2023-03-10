package mime

import (
	"fmt"
	"sort"
	"sync"

	"github.com/igorrendulic/couchdb-experiment/email/mime/parser"
)

var (
	parsersMu sync.RWMutex
	parsers   = make(map[string]parser.Parser)
)

// Register makes a parser available by the provided name.
// If Register is called twice with the same name or if parser is nil,
// it panics.
func Register(name string, parser parser.Parser) {
	parsersMu.Lock()
	defer parsersMu.Unlock()
	if parser == nil {
		panic("sql: Register driver is nil")
	}
	if _, dup := parsers[name]; dup {
		panic("sql: Register called twice for driver " + name)
	}
	parsers[name] = parser
}

func unregisterAllParsers() {
	parsersMu.Lock()
	defer parsersMu.Unlock()
	// For tests.
	parsers = make(map[string]parser.Parser)
}

// Parsers returns a sorted list of the names of the registered drivers.
func Parsers() []string {
	parsersMu.RLock()
	defer parsersMu.RUnlock()
	var list []string
	for name := range parsers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

// GetParser returns a parser by name.
func GetParser(parserName string) (parser.Parser, error) {
	parsersMu.RLock()
	defer parsersMu.RUnlock()
	parser, ok := parsers[parserName]
	if !ok {
		return nil, fmt.Errorf("sql: unknown parser %q (forgotten import?)", parserName)
	}
	return parser, nil
}
