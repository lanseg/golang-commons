package almostio

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"testing"
)

type FakeWriteCloser struct {
	io.WriteCloser

	Closed bool
	Data   []byte
}

func (fw *FakeWriteCloser) Write(b []byte) (int, error) {
	if fw.Data == nil {
		fw.Data = []byte{}
	}
	fw.Data = append(fw.Data, b...)
	return len(b), nil
}

func (fw *FakeWriteCloser) Close() error {
	fw.Closed = true
	return nil
}

func TestMultiWriteCloser(t *testing.T) {
	for _, tc := range []struct {
		name    string
		writers []*FakeWriteCloser
		closers []*FakeWriteCloser
		writes  [][]byte
		want    []byte
		wantErr bool
	}{
		{
			name: "no writers no writeclosers",
			writes: [][]byte{
				{0, 1, 2},
			},
		},
		{
			name:    "simple write",
			writers: []*FakeWriteCloser{{}, {}},
			closers: []*FakeWriteCloser{{}},
			writes: [][]byte{
				{0, 1, 2},
				{6, 5, 4},
			},
			want: []byte{0, 1, 2, 6, 5, 4},
		},
		{
			name:    "Empty write",
			writers: []*FakeWriteCloser{{}},
			closers: []*FakeWriteCloser{{}},
			writes:  [][]byte{{}, {}, {}, {1}},
			want:    []byte{1},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			closed := false
			mwc := NewMultiWriteCloser().SetOnClose(func() {
				closed = true
			})
			for _, w := range tc.writers {
				mwc.AddWriter(w)
			}
			for _, wc := range tc.closers {
				mwc.AddWriteCloser(wc)
			}
			for _, w := range tc.writes {
				mwc.Write(w)
			}
			mwc.Close()
			if !closed {
				t.Errorf("OnClose was not called during close")
			}
			for i, w := range append(tc.writers, tc.closers...) {
				if !reflect.DeepEqual(tc.want, w.Data) {
					t.Errorf("Expected rw %d to have data %v, but got %v", i, tc.want, w.Data)
				}
			}
			for i, w := range tc.writers {
				if w.Closed {
					t.Errorf("Writers should not be closed, but %d was", i)
				}
			}
			for i, w := range tc.closers {
				if !w.Closed {
					t.Errorf("WriterClosers should be closed, but %d was not", i)
				}
			}
		})
	}
}

func TestMultiWriteCloserRegressions(t *testing.T) {
	t.Run("Write returns incorrect byte count when using io.Copy", func(t *testing.T) {
		closed := false
		fwc1 := &FakeWriteCloser{}
		fwc2 := &FakeWriteCloser{}
		mwc := NewMultiWriteCloser().
			AddWriteCloser(fwc1).
			AddWriteCloser(fwc2).
			SetOnClose(func() {
				closed = true
			})

		toCopy := []byte("Whatever bytes")
		copied, err := io.Copy(mwc, bytes.NewBuffer(toCopy))

		if err != nil {
			t.Errorf("Error when copying bytes: %s", err)
			return
		}
		if copied != int64(len(toCopy)) {
			t.Errorf("Expected to copy %d bytes, but got %d.", len(toCopy), copied)
			return
		}
		fmt.Printf("HERE: %s %s\n", string(fwc1.Data), string(fwc2.Data))
		mwc.Close()
		if !closed {
			t.Errorf("OnClose was not called during close")
		}
	})
}

type sampleFile struct {
	originalFileName string
	content          []byte
}

func TestOverlay(t *testing.T) {

	for _, tc := range []struct {
		name          string
		originalFiles []*sampleFile
		expectedFiles []*sampleFile
	}{
		{
			name: "Single file name",
			originalFiles: []*sampleFile{{
				originalFileName: "Whatever",
				content:          []byte{1, 2, 3},
			}},
		},
		{
			name: "File name with spaces",
			originalFiles: []*sampleFile{{
				originalFileName: "A spaces file name",
				content:          []byte{1, 2, 3, 4, 5},
			}},
		},
		{
			name: "File name with unicode",
			originalFiles: []*sampleFile{{
				originalFileName: "Немного ユニコード",
				content:          []byte{1, 2, 3, 4, 5},
			}},
		},
		{
			name: "File name with path",
			originalFiles: []*sampleFile{{
				originalFileName: "some/file/path",
				content:          []byte{1, 2, 3, 4, 5},
			}},
		},
		{
			name: "Empty file",
			originalFiles: []*sampleFile{{
				originalFileName: "some empty file",
				content:          []byte{},
			}},
		},
		{
			name: "Two files different names",
			originalFiles: []*sampleFile{
				{"File 1", []byte{1, 2, 3}},
				{"File 2", []byte{4, 5, 6}},
			},
		},
		{
			name: "Two files same name",
			originalFiles: []*sampleFile{
				{"File 1", []byte{1, 2, 3}},
				{"File 1", []byte{4, 5, 6}},
			},
			expectedFiles: []*sampleFile{
				{"File 1", []byte{4, 5, 6}},
			},
		},
		{
			name: "PNG file header",
			originalFiles: []*sampleFile{
				{"File 1", []byte{
					0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
					0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
					0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
					0x01, 0x00, 0x00, 0x00, 0x00, 0x37, 0x6e, 0xf9,
					0x24, 0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41,
					0x54, 0x78, 0x01, 0x63, 0x60, 0x00, 0x00, 0x00,
					0x02, 0x00, 0x01, 0x73, 0x75, 0x01, 0x18, 0x00,
					0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae,
					0x42, 0x60, 0x82,
				}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			o, err := NewLocalOverlay(filepath.Join(t.TempDir(), "overlay_root"))
			if err != nil {
				t.Errorf("Error while starting an overlay: %v", err)
				return
			}
			for _, f := range tc.originalFiles {
				wc, err := o.OpenWrite("someid", f.originalFileName)
				if err != nil {
					t.Errorf("Error while writing file %v: %v", f, err)
					return
				}
				wc.Write(f.content)
				wc.Close()
			}
			if tc.expectedFiles == nil {
				tc.expectedFiles = tc.originalFiles
			}

			result := []*sampleFile{}
			for _, f := range tc.expectedFiles {
				reader, err := o.OpenRead("someid", f.originalFileName)
				content := []byte{}
				if err == nil {
					content, err = io.ReadAll(reader)
				}
				if err == nil {
					result = append(result, &sampleFile{f.originalFileName, content})
				}
			}

			if !reflect.DeepEqual(result, tc.expectedFiles) {
				t.Errorf("Expected result to be %v, but got %v", tc.expectedFiles, result)
			}
		})
	}
}
