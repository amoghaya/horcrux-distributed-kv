package server

//Connects parsed commands → KV engine

import (
	"fmt"
	"horcrux/internal/kv"
)

// Handler processes commands
type Handler struct {
	store *kv.KVStore
}

func NewHandler(store *kv.KVStore) *Handler {
	return &Handler{store: store}
}

// Execute runs command and returns response
func (h *Handler) Execute(cmd Command) string {
	switch cmd.Name {

	case "SET":
		if len(cmd.Args) < 2 {
			return "ERR wrong number of args"
		}
		h.store.Put(cmd.Args[0], cmd.Args[1])
		return "OK"

	case "GET":
		val, ok := h.store.Get(cmd.Args[0])
		if !ok {
			return "NULL"
		}
		return fmt.Sprintf("%v", val)

	case "DEL":
		h.store.Delete(cmd.Args[0])
		return "OK"

	default:
		return "ERR unknown command"
	}
}
