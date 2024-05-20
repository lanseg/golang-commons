package almostio

import (
	"io"
)

type fixedSizeWriter struct {
	io.Writer

	parent  io.Writer
	written int
	maxSize int
}

func (fsw *fixedSizeWriter) Write(b []byte) (int, error) {
	copySize := fsw.maxSize - fsw.written
	if copySize <= 0 {
		return 0, nil
	}
	if copySize > len(b) {
		copySize = len(b)
	}
	fsw.written += copySize
	_, err := fsw.parent.Write(b[:copySize])
	return len(b), err
}

func FixedSizeWriter(w io.Writer, maxCapacity int) io.Writer {
	return &fixedSizeWriter{
		parent:  w,
		written: 0,
		maxSize: maxCapacity,
	}
}
