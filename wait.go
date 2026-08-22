package ssehub

import "context"

func WaitFrame(ctx context.Context, ch <-chan []byte) ([]byte, error) {
	// 带截止时间的等帧：ctx 结束（超时/取消）后必须尽快收工，
	// 不能死堵在接收通道上耗着，否则会把整条网关链路的超时拖长。
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
