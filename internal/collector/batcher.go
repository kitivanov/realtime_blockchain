package collector

import (
	"realtime_blockchain/internal/model"
	"time"
)

type Batcher struct {
	in <-chan model.LogTrigger

	size     int
	interval time.Duration

	out chan []model.LogTrigger
}

func NewBatcher(in <-chan model.LogTrigger, size int, interval time.Duration) *Batcher {
	return &Batcher{
		in:       in,
		size:     size,
		interval: interval,
		out:      make(chan []model.LogTrigger, 100),
	}
}

func (b *Batcher) Out() <-chan []model.LogTrigger {
	return b.out
}

func (b *Batcher) Run() {
	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()

	var buffer []model.LogTrigger

	flush := func() {
		if len(buffer) == 0 {
			return
		}

		batch := make([]model.LogTrigger, len(buffer))
		copy(batch, buffer)

		b.out <- batch
		buffer = buffer[:0]
	}

	for {
		select {

		case v, ok := <-b.in:
			if !ok {
				flush()
				close(b.out)
				return
			}

			buffer = append(buffer, v)

			if len(buffer) >= b.size {
				flush()
			}

		case <-ticker.C:
			flush()
		}
	}
}
