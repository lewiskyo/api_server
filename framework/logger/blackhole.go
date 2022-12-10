package logger

import "io"

var _ io.WriteCloser = (*BlackHole)(nil)

type BlackHole struct{}

func (b BlackHole) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (b BlackHole) Close() error {
	return nil
}
