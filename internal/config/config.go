package config

import "time"

type Config struct {
	// Server
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxConnections int

	// Cluster
	ReplicationFactor int
	WriteQuorum       int
	ReadQuorum        int
	RingReplicas      int

	// Failure detection
	NodeTimeout   time.Duration
	HeartbeatFreq time.Duration

	// Storage engine
	ShardCount   int
	LRUSize      int
	TTLCheckFreq time.Duration

	// WAL
	WALPath     string
	SyncOnWrite bool
}

func DefaultConfig() *Config {
	return &Config{
		Port:           ":8080",
		ReadTimeout:    2 * time.Second,
		WriteTimeout:   2 * time.Second,
		MaxConnections: 1000,

		ReplicationFactor: 3,
		WriteQuorum:       2,
		ReadQuorum:        2,
		RingReplicas:      3,

		NodeTimeout:   2 * time.Second,
		HeartbeatFreq: 1 * time.Second,

		ShardCount:   16,
		LRUSize:      1000,
		TTLCheckFreq: 5 * time.Second,

		WALPath:     "wal.log",
		SyncOnWrite: true,
	}
}
