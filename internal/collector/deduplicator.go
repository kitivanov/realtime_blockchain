package collector

import (
	"realtime_blockchain/internal/model"
	"sync"
	"time"
)

type Deduplicator struct {
	in   <-chan model.LogTrigger
	out  chan model.LogTrigger
	seen map[string]int64
	mu   sync.Mutex
	ttl  time.Duration
}

func NewDeduplicator(in <-chan model.LogTrigger, ttl time.Duration) *Deduplicator {
	return &Deduplicator{
		in:   in,
		out:  make(chan model.LogTrigger, 1000),
		seen: make(map[string]int64),
		ttl:  ttl,
	}
}

func (d *Deduplicator) Out() <-chan model.LogTrigger {
	return d.out
}

func (d *Deduplicator) Seen(tx string) bool {
	now := time.Now().Unix()

	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.seen[tx]; ok {
		if now-t < int64(d.ttl.Seconds()) {
			return true
		}
	}

	d.seen[tx] = now
	return false
}

func (d *Deduplicator) Run() {
	for v := range d.in {
		if d.isDuplicate(v.TxHash) {
			continue
		}

		d.out <- v
	}
}

func (d *Deduplicator) isDuplicate(tx string) bool {
	now := time.Now().Unix()

	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.seen[tx]; ok {
		if now-t < int64(d.ttl.Seconds()) {
			return true
		}
	}

	d.seen[tx] = now
	return false
}

func (d *Deduplicator) StartCleanup() {
	ticker := time.NewTicker(time.Minute)

	go func() {
		for range ticker.C {
			now := time.Now().Unix()

			d.mu.Lock()
			for k, v := range d.seen {
				if now-v > int64(d.ttl.Seconds()) {
					delete(d.seen, k)
				}
			}
			d.mu.Unlock()
		}
	}()
}
