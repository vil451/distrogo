package cancelable_reader

import (
	"context"
	"io"
)

const MaxBufSize = 4096

type CancelableReader struct {
	ctx  context.Context
	data chan []byte
	err  error
	r    io.Reader
}

func (c *CancelableReader) Read() ([]byte, error) {
	select {
	case <-c.ctx.Done():
		return nil, c.ctx.Err()
	case d, ok := <-c.data:
		if !ok {
			return nil, c.err
		}
		return d, nil
	}
}

func New(ctx context.Context, r io.Reader) *CancelableReader {
	c := &CancelableReader{
		r:    r,
		ctx:  ctx,
		data: make(chan []byte),
	}
	go c.begin()
	return c
}

func (c *CancelableReader) begin() {
	buf := make([]byte, MaxBufSize)
	for {
		n, err := c.r.Read(buf)
		if n > 0 {
			tmp := make([]byte, n)
			copy(tmp, buf[:n])
			c.data <- tmp
		}
		if err != nil {
			c.err = err
			close(c.data)
			return
		}
	}
}
