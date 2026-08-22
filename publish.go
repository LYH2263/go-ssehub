package ssehub

import (
	"context"
	"fmt"
	"github.com/LYH2263/go-ssehub/internal/clone"
	"github.com/LYH2263/go-ssehub/internal/event"
)

func (r *Room) Publish(ctx context.Context, name string, data []byte) (event.Event, error) {
	if err := ctx.Err(); err != nil {
		return event.Event{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return event.Event{}, ErrClosed
	}
	if r.Enc == nil {
		return event.Event{}, fmt.Errorf("%w: encoder", ErrInvalid)
	}
	ev := event.Event{Name: name, Data: clone.Bytes(data)}
	ev = r.Replay.Append(ev)
	frame := r.Enc.Encode(ev)
	var first error
	for _, s := range r.Subs {
		if err := s.Q.Push(frame); err != nil {
			if first == nil {
				first = err
			}
		}
	}
	return ev, first
}
