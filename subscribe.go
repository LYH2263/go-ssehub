package ssehub

import (
	"context"
	"fmt"
	"github.com/LYH2263/go-ssehub/internal/backpressure"
)

func (r *Room) Subscribe(ctx context.Context, subID string, lastID int64, qCap int) (<-chan []byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return nil, ErrClosed
	}
	if r.Enc == nil {
		return nil, fmt.Errorf("%w: encoder", ErrInvalid)
	}
	if _, ok := r.Subs[subID]; ok {
		return nil, ErrConflict
	}
	q := backpressure.New(qCap)
	r.Subs[subID] = &Sub{ID: subID, Room: r.Name, Q: q}
	r.subOrder = append(r.subOrder, subID)
	for _, ev := range r.Replay.After(lastID) {
		frame := r.Enc.Encode(ev)
		_ = q.Push(frame)
	}
	return q.Chan(), nil
}
func (r *Room) Unsubscribe(subID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.Subs[subID]; ok {
		s.Q.Close()
		delete(r.Subs, subID)
		for i, id := range r.subOrder {
			if id == subID {
				r.subOrder = append(r.subOrder[:i], r.subOrder[i+1:]...)
				break
			}
		}
	}
}
