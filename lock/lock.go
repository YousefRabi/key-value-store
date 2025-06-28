package lock

import (
	"6.5840/kvtest1"
	"sync"
)

// Store locks using their key
var locks = make(map[string]*Lock)

type Lock struct {
	// IKVClerk is a go interface for k/v clerks: the interface hides
	// the specific Clerk type of ck but promises that ck supports
	// Put and Get.  The tester passes the clerk in when calling
	// MakeLock().
	ck  kvtest.IKVClerk
	key string
	mu  sync.Mutex
}

// The tester calls MakeLock() and passes in a k/v clerk; your code can
// perform a Put or Get by calling lk.ck.Put() or lk.ck.Get().
//
// Use l as the key to store the "lock state" (you would have to decide
// precisely what the lock state is).
func MakeLock(ck kvtest.IKVClerk, l string) *Lock {
	if lk, ok := locks[l]; ok {
		return lk
	}
	lk := &Lock{ck: ck, key: l}
	locks[l] = lk
	return lk
}

func (lk *Lock) Acquire() {
	lock := locks[lk.key]
	lock.mu.Lock()
}

func (lk *Lock) Release() {
	lock := locks[lk.key]
	lock.mu.Unlock()
}
