package almostio

import (
	"bytes"
	"io"
	"math/rand"
	"os"
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

func makeString(len int) string {
	result := make([]byte, len)
	for i := range len {
		result[i] = byte(rand.Intn('Z'-'A') + 'A')
	}
	return string(result)
}

func TestOverlay(t *testing.T) {

	t.Run("Create root folder if none", func(t *testing.T) {
		tmp := t.TempDir()
		NewLocalOverlay(filepath.Join(tmp, "overlay_root"), NewJsonMarshal[OverlayMetadata]())
		if _, err := os.Stat(filepath.Join(tmp, "overlay_root", systemFolderName, metadataFileName)); os.IsNotExist(err) {
			t.Errorf("Overlay structure not created")
		}
	})

	t.Run("Create empty folder if no system folder", func(t *testing.T) {
		tmp := t.TempDir()
		os.MkdirAll(filepath.Join(tmp, "overlay_root"), defaultDirPermissions)
		NewLocalOverlay(filepath.Join(tmp, "overlay_root"), NewJsonMarshal[OverlayMetadata]())
		if _, err := os.Stat(filepath.Join(tmp, "overlay_root", systemFolderName, metadataFileName)); os.IsNotExist(err) {
			t.Errorf("Overlay structure not created")
		}
	})

	t.Run("Create empty file if no metadata file", func(t *testing.T) {
		tmp := t.TempDir()
		os.MkdirAll(filepath.Join(tmp, "overlay_root", systemFolderName), defaultDirPermissions)
		NewLocalOverlay(filepath.Join(tmp, "overlay_root"), NewJsonMarshal[OverlayMetadata]())
		if _, err := os.Stat(filepath.Join(tmp, "overlay_root", systemFolderName, metadataFileName)); os.IsNotExist(err) {
			t.Errorf("Overlay structure not created")
		}
	})

	t.Run("Create empty file if no metadata file", func(t *testing.T) {
		tmp := t.TempDir()
		os.MkdirAll(filepath.Join(tmp, "overlay_root", systemFolderName), defaultDirPermissions)
		if _, err := NewLocalOverlay(filepath.Join(tmp, "overlay_root"), NewJsonMarshal[OverlayMetadata]()); err != nil {
			t.Errorf("Cannot create overlay: %s", err)
		}
		if _, err := os.Stat(filepath.Join(tmp, "overlay_root", systemFolderName, metadataFileName)); os.IsNotExist(err) {
			t.Errorf("Overlay structure not created")
		}
	})

	veryLongName := makeString(1000)
	for _, tc := range []struct {
		name             string
		originalFiles    []*sampleFile
		expectedFiles    []*sampleFile
		expectedMetadata []*FileMetadata
	}{
		{
			name: "Single file name",
			originalFiles: []*sampleFile{{
				originalFileName: "Whatever",
				content:          []byte{1, 2, 3},
			}},
			expectedMetadata: []*FileMetadata{{
				Name:      "Whatever",
				LocalName: "d4ae358f_Whatever",
				Sha256:    "039058c6f2c0cb492c533b0a4d14ef77cc0f78abccced5287d84a1a2011cfb81",
				Mime:      "application/octet-stream",
			}},
		},
		{
			name: "File name with spaces",
			originalFiles: []*sampleFile{{
				originalFileName: "A spaces file name",
				content:          []byte{1, 2, 3, 4, 5},
			}},
			expectedMetadata: []*FileMetadata{{
				Name:      "A spaces file name",
				LocalName: "89493756_A_spaces_file_name",
				Sha256:    "74f81fe167d99b4cb41d6d0ccda82278caee9f3e2f25d5e5a3936ff3dcec60d0",
				Mime:      "application/octet-stream",
			}},
		},
		{
			name: "File name with unicode",
			originalFiles: []*sampleFile{{
				originalFileName: "Немного ユニコード",
				content:          []byte{1, 2, 3, 4, 5},
			}},
			expectedMetadata: []*FileMetadata{{
				Name:      "Немного ユニコード",
				LocalName: "7ab716af__",
				Sha256:    "74f81fe167d99b4cb41d6d0ccda82278caee9f3e2f25d5e5a3936ff3dcec60d0",
				Mime:      "application/octet-stream",
			}},
		},
		{
			name: "File name with path",
			originalFiles: []*sampleFile{{
				originalFileName: "some/file/path",
				content:          []byte{1, 2, 3, 4, 5},
			}},
			expectedMetadata: []*FileMetadata{{
				Name:      "some/file/path",
				LocalName: "d7aeb9d2_some_file_path",
				Sha256:    "74f81fe167d99b4cb41d6d0ccda82278caee9f3e2f25d5e5a3936ff3dcec60d0",
				Mime:      "application/octet-stream",
			}},
		},
		{
			name: "Empty file",
			originalFiles: []*sampleFile{{
				originalFileName: "some empty file",
				content:          []byte{},
			}},
			expectedMetadata: []*FileMetadata{{
				Name:      "some empty file",
				LocalName: "d78d0606_some_empty_file",
				Sha256:    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				Mime:      "text/plain; charset=utf-8",
			}},
		},
		{
			name: "Two files different names",
			originalFiles: []*sampleFile{
				{"File 1", []byte{1, 2, 3}},
				{"File 2", []byte{4, 5, 6}},
			},
			expectedMetadata: []*FileMetadata{
				{
					Name:      "File 1",
					LocalName: "644fe258_File_1",
					Sha256:    "039058c6f2c0cb492c533b0a4d14ef77cc0f78abccced5287d84a1a2011cfb81",
					Mime:      "application/octet-stream",
				},
				{
					Name:      "File 2",
					LocalName: "674fe711_File_2",
					Sha256:    "787c798e39a5bc1910355bae6d0cd87a36b2e10fd0202a83e3bb6b005da83472",
					Mime:      "application/octet-stream",
				},
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
			expectedMetadata: []*FileMetadata{{
				Name:      "File 1",
				LocalName: "644fe258_File_1",
				Sha256:    "787c798e39a5bc1910355bae6d0cd87a36b2e10fd0202a83e3bb6b005da83472",
				Mime:      "application/octet-stream",
			}},
		},
		{
			name: "Very long file name",
			originalFiles: []*sampleFile{
				{veryLongName, []byte{1, 2, 3}},
				{veryLongName + "1", []byte{4, 5, 6}},
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
			expectedMetadata: []*FileMetadata{{
				Name:      "File 1",
				LocalName: "644fe258_File_1",
				Sha256:    "836c5e8c94b74be78456122528bb2a44b4cf61b3922f211e2bae5bf327f95f09",
				Mime:      "image/png",
			}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			o, err := NewLocalOverlay(filepath.Join(t.TempDir(), "overlay_root"), NewJsonMarshal[OverlayMetadata]())
			if err != nil {
				t.Errorf("Error while starting an overlay: %v", err)
				return
			}
			for _, f := range tc.originalFiles {
				wc, err := o.OpenWrite(f.originalFileName)
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
				reader, err := o.OpenRead(f.originalFileName)
				if err != nil {
					t.Errorf("Error while opening %s for reading: %s", f.originalFileName, err)
					continue
				}
				content, err := io.ReadAll(reader)
				if err != nil {
					t.Errorf("Error while reading from  %s: %s", f.originalFileName, err)
					continue
				}
				result = append(result, &sampleFile{f.originalFileName, content})
			}

			if !reflect.DeepEqual(result, tc.expectedFiles) {
				t.Errorf("Expected result to be %v, but got %v", tc.expectedFiles, result)
			}

			if tc.expectedMetadata != nil {
				names := []string{}
				for _, fn := range tc.expectedFiles {
					names = append(names, fn.originalFileName)
				}

				md := o.GetMetadata(names)
				if len(md) != len(tc.expectedMetadata) {
					t.Errorf("Expected metadata to have length %d, but got %d", len(md), len(tc.expectedMetadata))
				}
				for i, mdata := range md {
					if !reflect.DeepEqual(tc.expectedMetadata[i], mdata) {
						t.Errorf("Metadata for %s differs", names[i])
					}
				}
			}
		})
	}
}
