package ssehub

import "context"

func WaitFrame(ctx context.Context, ch <-chan []byte) ([]byte, error) {
	_ = ctx
	b, ok := <-ch
	if !ok {
		return nil, ErrClosed
	}
	return b, nil
}
