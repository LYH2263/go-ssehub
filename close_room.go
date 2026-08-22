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
	// Mark the room closed with an explicit flag. The closed *state* must be
	// decided by r.closed, never inferred from r.Replay being nil/empty:
	// Publish checks r.closed before touching r.Replay, so nil-ing Replay here
	// is a cleanup, not the close signal. (See publish.go and the
	// publish_close_test.go regression.)
	r.Replay = nil
	r.closed = true
}

func (r *Room) CloseFlushCount() int {
	r.Close()
	return 0
}
