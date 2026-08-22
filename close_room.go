package ssehub

func (r *Room) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return
	}
	r.Replay.Clear()
	for id, s := range r.Subs {
		s.Q.Close()
		delete(r.Subs, id)
	}
	r.subOrder = nil
	r.closed = true
}
func (r *Room) CloseFlushCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return 0
	}
	r.Replay.Clear()
	flushed := r.Replay.Flush()
	for id, s := range r.Subs {
		s.Q.Close()
		delete(r.Subs, id)
	}
	r.subOrder = nil
	r.closed = true
	return len(flushed)
}
