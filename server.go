package kvsrv

import (
	"log"
	"sync"

	"6.5840/labrpc"
	"6.5840/rpc"
	"6.5840/tester1"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type ValueVersion struct {
	Value   string
	Version rpc.Tversion
}

type KVServer struct {
	mu sync.Mutex

	kvMap map[string]*ValueVersion
}

func MakeKVServer() *KVServer {
	kv := &KVServer{}
	kv.kvMap = make(map[string]*ValueVersion)
	return kv
}

// Get returns the value and version for args.Key, if args.Key
// exists. Otherwise, Get returns ErrNoKey.
func (kv *KVServer) Get(args *rpc.GetArgs, reply *rpc.GetReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	valueVersion, ok := kv.kvMap[args.Key]

	if !ok {
		reply.Err = rpc.ErrNoKey
		return
	}

	reply.Value = valueVersion.Value
	reply.Version = valueVersion.Version
	reply.Err = rpc.OK
}

// Update the value for a key if args.Version matches the version of
// the key on the server. If versions don't match, return ErrVersion.
// If the key doesn't exist, Put installs the value if the
// args.Version is 0, and returns ErrNoKey otherwise.
func (kv *KVServer) Put(args *rpc.PutArgs, reply *rpc.PutReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	valueVersion, ok := kv.kvMap[args.Key]

	if !ok && args.Version == 0 {
		kv.kvMap[args.Key] = &ValueVersion{Value: args.Value, Version: 1}
		reply.Err = rpc.OK
		return
	}

	if !ok {
		reply.Err = rpc.ErrNoKey
		return
	}

	if args.Version != valueVersion.Version {
		reply.Err = rpc.ErrVersion
		return
	}

	valueVersion.Value = args.Value
	valueVersion.Version = args.Version + 1
	reply.Err = rpc.OK

}

// You can ignore Kill() for this lab
func (kv *KVServer) Kill() {
}

// You can ignore all arguments; they are for replicated KVservers
func StartKVServer(ends []*labrpc.ClientEnd, gid tester.Tgid, srv int, persister *tester.Persister) []tester.IService {
	kv := MakeKVServer()
	return []tester.IService{kv}
}
