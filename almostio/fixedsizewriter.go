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

// Write sends bytes to the underlying writer, but not more than given amount.
// Always returns len(b) as the result to avoid "short write" errors
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

// FixedSizeWriter creates writer that forwards at most maxCapacity byte to the underlying writer.
func FixedSizeWriter(w io.Writer, maxCapacity int) io.Writer {
	return &fixedSizeWriter{
		parent:  w,
		written: 0,
		maxSize: maxCapacity,
	}
}
