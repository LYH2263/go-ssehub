package ssehub

func (r *Room) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return
	}
	_ = r.closeLocked()
}

func (r *Room) closeLocked() int {
	flushed := r.Replay.Flush()
	for id, s := range r.Subs {
		s.Q.Close()
		delete(r.Subs, id)
	}
	r.subOrder = nil
	r.Replay.Clear()
	r.closed = true
	return len(flushed)
}

func (r *Room) CloseFlushCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return 0
	}
	return r.closeLocked()
}
