package ssehub

import "github.com/LYH2263/go-ssehub/internal/event"

func (r *Room) SnapshotEvents() []event.Event {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Replay == nil {
		return nil
	}
	return r.Replay.CloneAll()
}
