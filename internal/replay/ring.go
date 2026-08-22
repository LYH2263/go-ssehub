package replay

import (
	"github.com/LYH2263/go-ssehub/internal/clone"
	"github.com/LYH2263/go-ssehub/internal/event"
	"sync"
)

type Ring struct {
	mu   sync.Mutex
	buf  []event.Event
	cap  int
	next int64
}

func New(cap int) *Ring {
	if cap <= 0 {
		cap = 64
	}
	return &Ring{buf: make([]event.Event, 0, cap), cap: cap}
}
func (r *Ring) Append(ev event.Event) event.Event {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	ev.ID = r.next
	// Hold an independent copy of Data so later mutations of the caller's
	// buffer (reused across frames) cannot retroactively dirty the ring.
	stored := ev
	stored.Data = clone.Bytes(ev.Data)
	if len(r.buf) < r.cap {
		r.buf = append(r.buf, stored)
	} else {
		r.buf = append(r.buf[1:], stored)
	}
	return ev
}
func (r *Ring) After(id int64) []event.Event {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]event.Event, 0)
	for _, ev := range r.buf {
		if ev.ID > id {
			cp := ev
			cp.Data = append([]byte(nil), ev.Data...)
			out = append(out, cp)
		}
	}
	return out
}
func (r *Ring) Flush() []event.Event {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]event.Event, len(r.buf))
	for i, ev := range r.buf {
		out[i] = ev
		out[i].Data = append([]byte(nil), ev.Data...)
	}
	return out
}
func (r *Ring) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.buf = r.buf[:0]
}

func (r *Ring) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.buf)
}

func (r *Ring) CloneAll() []event.Event {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]event.Event, len(r.buf))
	for i, ev := range r.buf {
		out[i] = ev
		out[i].Data = append([]byte(nil), ev.Data...)
	}
	return out
}
