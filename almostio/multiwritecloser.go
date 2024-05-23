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

// MultiWriteCloser is similar to the io.MultiWriter, but with the Close() method.
type MultiWriteCloser struct {
	io.WriteCloser

	onClose func()
	writers []io.WriteCloser
}

// NewMultiWriteCloser creates default MultiWriteCloser with empty writers and no onClose function.
func NewMultiWriteCloser() *MultiWriteCloser {
	return &MultiWriteCloser{
		writers: []io.WriteCloser{},
	}
}

// AddWriter adds another writer that will accept data written to the MultiWriteCloser.
func (mw *MultiWriteCloser) AddWriter(w io.Writer) *MultiWriteCloser {
	mw.writers = append(mw.writers, NopWriteCloser(w))
	return mw
}

// AddWriterCloser adds another writerCloser that will accept data written to the MultiWriteCloser.
// Its Close() method invoked when parent onCall is invoked.
func (mw *MultiWriteCloser) AddWriteCloser(w io.WriteCloser) *MultiWriteCloser {
	mw.writers = append(mw.writers, w)
	return mw
}

// SetOnClose configures a function that is called after all child WriteClosers successfuly closed.
func (mw *MultiWriteCloser) SetOnClose(onClose func()) *MultiWriteCloser {
	mw.onClose = onClose
	return mw
}

// Close invokes "Close" for all underlying WriteClosers, returns an error if any of them fails.
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

// Write writes given bytes to all the underlying writers.
func (mw *MultiWriteCloser) Write(b []byte) (int, error) {
	for _, wc := range mw.writers {
		written, err := wc.Write(b)
		if err != nil {
			return 0, err
		}
		if written != len(b) {
			return 0, io.ErrShortWrite
		}
	}
	return len(b), nil
}
