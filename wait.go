package ssehub

import "context"

func WaitFrame(ctx context.Context, ch <-chan []byte) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case b, ok := <-ch:
		if !ok {
			return nil, ErrClosed
		}
		return b, nil
	}
}
