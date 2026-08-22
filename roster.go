package ssehub

import "strings"

func (r *Room) ExportRoster() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return strings.Join(r.subOrder, ",")
}
func (r *Room) HasSubscriber(subID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.Subs[subID]
	return ok
}
