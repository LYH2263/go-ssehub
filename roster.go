package ssehub

import "strings"

func (r *Room) ExportRoster() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	ids := make([]string, len(r.subOrder))
	copy(ids, r.subOrder)
	return strings.Join(ids, ",")
}

func (r *Room) HasSubscriber(subID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.Subs[subID]
	return ok
}
