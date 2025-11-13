package window

import "sync"

type Rolling struct {
	data   []float64
	size   int
	idx    int
	filled bool
	mu     sync.RWMutex
}

func New(size int) *Rolling {
	return &Rolling{
		data: make([]float64, size),
		size: size,
	}
}

func (r *Rolling) Push(v float64) {
	r.mu.Lock()
	r.data[r.idx] = v
	r.idx++
	if r.idx == r.size {
		r.idx = 0
		r.filled = true
	}
	r.mu.Unlock()
}

func (r *Rolling) Snapshot() []float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.filled {
		out := make([]float64, r.idx)
		copy(out, r.data[:r.idx])
		return out
	}

	out := make([]float64, r.size)
	n := r.size - r.idx
	copy(out, r.data[r.idx:])
	copy(out[n:], r.data[:r.idx])
	return out
}
