package coredns_blackhole

import (
	"sync"
	"time"
)

var BlMu sync.RWMutex

type Blocklist struct {
	list map[string]struct{}
}

func NewBlocklist() *Blocklist {
	blocklist := &Blocklist{}
	blocklist.list = map[string]struct{}{}
	return blocklist
}

func (b Blocklist) Add(s string) {
	BlMu.Lock()
	defer BlMu.Unlock()
	b.list[s] = struct{}{}
}

func (b Blocklist) Find(s string) bool {
	BlMu.RLock()
	defer BlMu.RUnlock()
	_, ok := b.list[s]
	return ok
}

func (b Blocklist) Len() (length int) {
	BlMu.RLock()
	defer BlMu.RUnlock()
	length = len(b.list)
	return length
}

type Options struct {
	refresh time.Duration
	retry   int
}

func NewOptions() Options {
	o := Options{}
	o.refresh = 30
	o.retry = 3
	return o
}
