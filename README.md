![Go](https://img.shields.io/badge/Go-1.25-blue)
![Status](https://img.shields.io/badge/status-learning_project-green)
![License](https://img.shields.io/badge/license-MIT-purple)


## Horcrux

Horcrux is a production-inspired *fault tolerant - distributed key-value storage* engine built in Go.

It started as a simple in-memory key-value store slowly evolved into a deeper systems engineering project exploring concepts used in real distributed databases like Redis, DynamoDB, and Cassandra.

The goal of this project was not just to “build a database”, but to understand how modern backend systems handle durability, concurrency, replication, fault tolerance, and performance under load.

---

## Features

### Core KV Operations
- PUT
- GET
- DELETE
- TTL-based expiration

### Storage Engine
- In-memory storage engine
- Modular storage abstraction
- Pluggable architecture for future storage engines

### Durability
- Write-Ahead Logging (WAL)
- Crash recovery using WAL replay

### Concurrency
- Thread-safe operations
- Lock striping / shard locking
- Race-condition tested using Go race detector

### Eviction
- LRU eviction policy
- Configurable memory capacity

### Distributed Systems Concepts
- Consistent hashing
- Replication
- Read quorum / write quorum
- Read repair
- Failure detection
- Rebalancing on node addition

### Performance & Observability
- Latency metrics
- Allocation benchmarking
- Object pooling using `sync.Pool`

---


## Architecture

```text
                    ┌────────────────────┐
                    │      Client        │
                    │  (telnet / CLI)    │
                    └─────────┬──────────┘
                              │ TCP (raw text protocol)
                              ▼
                    ┌────────────────────┐
                    │    TCP Server      │
                    └─────────┬──────────┘
                              │
                              ▼
                    ┌────────────────────┐
                    │     Handler        │
                    │ Command Processor  │
                    └─────────┬──────────┘
                              │
                              ▼
                    ┌────────────────────┐
                    │   Coordinator      │
                    │ Routing + Quorum   │
                    └─────────┬──────────┘
                              │
        ┌──────────────────────┼──────────────────────┐
        │                      │                      │
        ▼                      ▼                      ▼
┌────────────────┐   ┌────────────────┐   ┌────────────────┐
│    Node A      │   │    Node B      │   │    Node C      │
│ (Replica #1)   │   │ (Replica #2)   │   │ (Replica #3)   │
└───────┬────────┘   └───────┬────────┘   └───────┬────────┘
        │                    │                    │
        ▼                    ▼                    ▼
┌───────────────────────────────────────────────────────┐
│            InMemory KV Engine (per node)              │
│───────────────────────────────────────────────────────│
│ • HashMap storage                                     │
│ • Sharded locks (concurrency control)                 │
│ • TTL handling                                        │
│ • LRU eviction policy                                 │
│ • WAL (write-ahead log)                               │
│ • Metrics / latency tracking                          │
└───────────────────────────────────────────────────────┘

```
> Note:
> Distributed node communication is currently simulated in-process.
> Nodes are represented as independent storage instances running inside a single Go runtime.

---

## Project Structure

```text
internal/
│
├── kv/                # Storage engine, WAL, eviction, metrics
├── cluster/           # Replication, hashing, coordinator
├── server/            # HTTP server and handlers
├── tests/             # Unit tests and benchmarks
```

---

## Running the Project

#### Clone the Repository

```bash
git clone https://github.com/YOUR_USERNAME/horcrux.git

cd horcrux
```

---

#### Run with Docker

Build image:

```bash
docker build -t horcrux .
```

Run container:

```bash
docker run -p 8080:8080 horcrux
```

---

#### Start CLI Client

In another terminal:

```bash
go run ./cmd/cli
```

---

## Example Commands

```text
SET name horcrux
GET name
DEL name
```

<img width="763" height="313" alt="Screenshot 2026-05-13 001116" src="https://github.com/user-attachments/assets/c4c93209-d646-4c79-8a87-d4ab7109931b" />


---



## Testing

Run all tests:

```bash
go test ./internal/kv -v
```

Run race detector:

```bash
go test ./internal/kv -race
```

Run benchmarks:

```bash
go test ./internal/kv -bench=. -benchmem
```

---

## Benchmark Results

<img width="917" height="227" alt="Screenshot 2026-05-13 011922" src="https://github.com/user-attachments/assets/3578aeec-11ea-401c-b8f6-7f3ba76d991c" />



## Future Improvements

- Replace in-memory nodes with networked gRPC-based replicas
- Add gossip-based failure detection (SWIM protocol)
- Introduce LSM-tree based disk persistence layer
- Improve read repair using version vectors instead of simple overwrite

---

## Inspiration

This project was heavily inspired by systems like:
- Redis
- DynamoDB
- Cassandra
- Riak

---

## Author

Built as a systems engineering and distributed systems learning project using Go.
