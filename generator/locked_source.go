package generator

import (
	"math/rand"
	"sync"
)

// copy from go lib, to make the rand thread-safe

// Source64 source of rand
type Source64 interface {
	rand.Source
	Uint64() uint64
}

// LockedSource source with a lock
type LockedSource struct {
	lk  sync.Mutex
	src Source64
}

func (r *LockedSource) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
}

func (r *LockedSource) Uint64() (n uint64) {
	r.lk.Lock()
	n = r.src.Uint64()
	r.lk.Unlock()
	return
}

func (r *LockedSource) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}

// seedPos implements Seed for a LockedSource without a race condition.
func (r *LockedSource) seedPos(seed int64, readPos *int8) {
	r.lk.Lock()
	r.src.Seed(seed)
	*readPos = 0
	r.lk.Unlock()
}

// read implements Read for a LockedSource without a race condition.
func (r *LockedSource) read(p []byte, readVal *int64, readPos *int8) (n int, err error) {
	r.lk.Lock()
	n, err = read(p, r.src.Int63, readVal, readPos)
	r.lk.Unlock()
	return
}

func read(p []byte, int63 func() int64, readVal *int64, readPos *int8) (n int, err error) {
	pos := *readPos
	val := *readVal
	for n = 0; n < len(p); n++ {
		if pos == 0 {
			val = int63()
			pos = 7
		}
		p[n] = byte(val)
		val >>= 8
		pos--
	}
	*readPos = pos
	*readVal = val
	return
}
