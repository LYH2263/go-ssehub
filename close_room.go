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
	// 先刷出真实待重放量，再清空环；顺序颠倒会把计数抹成 0，末班对账失真。
	flushed := r.Replay.Flush()
	r.Replay.Clear()
	for id, s := range r.Subs {
		s.Q.Close()
		delete(r.Subs, id)
	}
	r.subOrder = nil
	r.closed = true
	return len(flushed)
}
