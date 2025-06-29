package lock

import (
	"6.5840/kvtest1"
	"6.5840/rpc"
	"time"
)

type Lock struct {
	// IKVClerk is a go interface for k/v clerks: the interface hides
	// the specific Clerk type of ck but promises that ck supports
	// Put and Get.  The tester passes the clerk in when calling
	// MakeLock().
	ck       kvtest.IKVClerk
	key      string
	clientId string
}

// The tester calls MakeLock() and passes in a k/v clerk; your code can
// perform a Put or Get by calling lk.ck.Put() or lk.ck.Get().
//
// Use l as the key to store the "lock state" (you would have to decide
// precisely what the lock state is).
func MakeLock(ck kvtest.IKVClerk, l string) *Lock {
	lk := &Lock{ck: ck}
	lk.key = l
	lk.clientId = kvtest.RandValue(8)

	return lk
}

// Acquire the lock from the distributed key-value store.
// If kv[lk.key] == ”, that means it's available to be acquired by any client.
// Otherwise, acquire will block and wait for the lock to be available.
// If the put response for acquiring the key returns an ErrMaybe,
// that means either the response got lost and the lock might have or might not have been
// acquired by the client successfully.
// We do a get request to check the state of the lock and if the current client has acquired
// it, everything's okay. Otherwise, we retry the put as long as the lock is available.
func (lk *Lock) Acquire() {
	tries := 0
	var err rpc.Err
	for {
		if tries >= 1 {
			time.Sleep(10 * time.Millisecond)
		}
		tries++
		clientId, version, _ := lk.ck.Get(lk.key)

		if clientId == lk.clientId {
			return
		}

		if clientId == "" {
			err = lk.ck.Put(lk.key, lk.clientId, version)
			if err == rpc.OK {
				return
			}
		}
	}
}

// Release is similar to Acquire above.
// Client gets the lock first. If they are the holder, they update it to
// empty string with the correct version.
// If the response was ErrMaybe, they try again as above.
// If response was OK, all done.
func (lk *Lock) Release() {
	tries := 0
	var err rpc.Err
	for {
		if tries >= 1 {
			time.Sleep(10 * time.Millisecond)
		}
		tries++
		clientId, version, _ := lk.ck.Get(lk.key)

		if clientId != lk.clientId {
			return
		}

		err = lk.ck.Put(lk.key, "", version)
		if err == rpc.OK {
			return
		}
	}
}
