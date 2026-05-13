package main

import (
	"horcrux/internal/kv"
	"horcrux/internal/server"
)

func main() {
	wal, err := kv.NewWAL("wal.log")
	if err != nil {
		panic(err)
	}
	defer wal.Close()

	lru := kv.NewLRUPolicy(100)

	engine := kv.NewInMemoryEngine(lru, wal)

	//  RECOVERY STEP
	wal.Replay(engine)

	store := kv.NewKVStore(engine)

	handler := server.NewHandler(store)
	srv := server.NewServer(":8080", handler)

	srv.Start()
}
