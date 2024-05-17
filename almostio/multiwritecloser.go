package almostio

import (
	"io"
)

type nopWriteCloser struct {
	io.WriteCloser

	onClose func()
	writer  io.Writer
}

func NopWriteCloser(w io.Writer) io.WriteCloser {
	return &nopWriteCloser{
		writer: w,
	}
}

func (nop *nopWriteCloser) Close() error {
	if nop.onClose != nil {
		nop.onClose()
	}
	return nil
}

func (nop *nopWriteCloser) Write(p []byte) (int, error) {
	return nop.writer.Write(p)
}

type MultiWriteCloser struct {
	io.WriteCloser

	onClose func()
	writers []io.WriteCloser
}

func NewMultiWriteCloser() *MultiWriteCloser {
	return &MultiWriteCloser{
		writers: []io.WriteCloser{},
	}
}

func (mw *MultiWriteCloser) AddWriter(w io.Writer) *MultiWriteCloser {
	mw.writers = append(mw.writers, NopWriteCloser(w))
	return mw
}

func (mw *MultiWriteCloser) AddWriteCloser(w io.WriteCloser) *MultiWriteCloser {
	mw.writers = append(mw.writers, w)
	return mw
}

func (mw *MultiWriteCloser) SetOnClose(onClose func()) *MultiWriteCloser {
	mw.onClose = onClose
	return mw
}

func (mw *MultiWriteCloser) Close() error {
	for _, wc := range mw.writers {
		if err := wc.Close(); err != nil {
			return err
		}
	}
	if mw.onClose != nil {
		mw.onClose()
	}
	return nil
}

func (mw *MultiWriteCloser) Write(b []byte) (int, error) {
	for _, wc := range mw.writers {
		written, err := wc.Write(b)
		if written != len(b) {
			return 0, io.ErrShortWrite
		}
		if err != nil {
			return 0, err
		}
	}
	return len(b), nil
}
