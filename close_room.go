package ssehub

func (r *Room) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return
	}
	for id, s := range r.Subs {
		s.Q.Close()
		delete(r.Subs, id)
	}
	r.subOrder = nil
	r.Replay = nil
	r.closed = true
}

func (r *Room) CloseFlushCount() int {
	r.Close()
	return 0
}
