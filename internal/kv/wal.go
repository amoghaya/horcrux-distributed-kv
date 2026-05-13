package kv

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
)

// WAL logs write operations before applying them.
type WAL struct {
	file *os.File
	mu   sync.Mutex
}

// NewWAL opens/creates WAL file
func NewWAL(path string) (*WAL, error) {
	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_APPEND|os.O_RDWR,
		0644,
	)

	if err != nil {
		return nil, err
	}

	return &WAL{
		file: file,
	}, nil
}

// LogWrite appends operation to WAL
func (w *WAL) LogWrite(
	op string,
	key string,
	value interface{},
) {
	w.mu.Lock()
	defer w.mu.Unlock()

	var buf bytes.Buffer

	buf.WriteString(op)
	buf.WriteByte('|')
	buf.WriteString(key)
	buf.WriteByte('|')
	buf.WriteString(fmt.Sprintf("%v", value))
	buf.WriteByte('\n')

	w.file.Write(buf.Bytes())

	// durability guarantee
	w.file.Sync()
}

// Replay rebuilds engine state from WAL
func (w *WAL) Replay(engine *InMemoryEngine) error {
	file, err := os.Open(w.file.Name())
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, "|")

		// malformed entry protection
		if len(parts) < 2 {
			continue
		}

		op := parts[0]
		key := parts[1]

		switch op {

		case "PUT":
			if len(parts) < 3 {
				continue
			}

			value := parts[2]
			engine.store[key] = NewEntry(value)

		case "PUT_TTL":
			if len(parts) < 3 {
				continue
			}

			value := parts[2]
			engine.store[key] = NewEntry(value)

		case "DEL":
			delete(engine.store, key)
		}
	}

	return scanner.Err()
}

// Close WAL file
func (w *WAL) Close() {
	w.file.Close()
}
