# Key-Value Store with Distributed Locking

A linearizable key-value store implementation for MIT 6.5840 Distributed Systems Course (Lab 2), featuring at-most-once semantics and distributed lock coordination.

## Overview

This project implements a fault-tolerant key-value server that maintains linearizability and handles network failures gracefully. The system supports versioned key-value operations and provides a distributed locking mechanism built on top of the key-value store.

The original lab handout can be found [here](https://pdos.csail.mit.edu/6.824/labs/lab-kvsrv1.html).

## Features

- **Linearizable Operations**: All operations appear to execute atomically in some sequential order
- **At-Most-Once Semantics**: Put operations execute at most once despite network failures and retransmissions
- **Version-Based Updates**: Conditional updates using version numbers to prevent conflicts
- **Network Fault Tolerance**: Automatic retry mechanisms for dropped RPC requests/replies
- **Distributed Locking**: Lock coordination across multiple clients using the key-value store

## Architecture

### Core Components

- **KVServer**: The main key-value server maintaining an in-memory map with version control
- **Clerk**: Client-side interface for interacting with the server via RPC
- **Lock**: Distributed lock implementation using the key-value store for coordination

### Key Operations

#### Put(key, value, version)
- Installs or replaces a value for a key only if version numbers match
- Creates new keys when version is 0
- Returns `ErrVersion` for version mismatches, `ErrNoKey` for non-existent keys
- Handles `ErrMaybe` responses for ambiguous retry scenarios

#### Get(key)
- Retrieves the current value and version for a key
- Returns `ErrNoKey` if the key doesn't exist
- Safe for retransmission (read-only operation)

## Implementation Details

### Version Control System
Each key maintains a `(value, version)` tuple where:
- Version starts at 1 for new keys
- Version increments on successful updates
- Version mismatches prevent conflicting updates

### Network Reliability
- Client automatically retries failed RPCs
- Distinguishes between definitive failures and ambiguous cases
- Returns `ErrMaybe` when Put operations may have succeeded but responses were lost

### Thread Safety
- Server uses mutex locks to ensure concurrent access safety
- All operations are atomic from the client's perspective

## Usage

### Basic Key-Value Operations

```go
clerk := MakeClerk(client, serverName)

err := clerk.Put("mykey", "myvalue", 0)

value, version, err := clerk.Get("mykey")

err = clerk.Put("mykey", "newvalue", version)
```

### Distributed Locking

The lock implementation uses the key-value store for coordination across multiple clients. Each lock client gets a unique ID and stores lock ownership state in the KV store.

```go
lock := MakeLock(clerk, "lockname")

lock.Acquire()

// Critical section
// ... protected operations ...

lock.Release()
```

## File Structure

```
├── client.go          # Client-side RPC interface
├── server.go          # Key-value server implementation
├── rpc/
│   └── rpc.go         # RPC message definitions
├── lock/
│   └── lock.go        # Distributed lock implementation
├── kvtest1/           # Testing utilities
├── labrpc/            # RPC framework
└── tester1/           # Test infrastructure
```

## Build and Run

```bash
# Install dependencies
go mod download

# Run all tests
go test -v

# Run specific test suites
go test -v -run Reliable
go test -v -race  # Check for race conditions

# Test locks
cd lock && go test -v
```
