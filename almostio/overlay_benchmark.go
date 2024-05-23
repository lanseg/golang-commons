package almostio

import (
	"fmt"
	"path/filepath"
	"testing"
)

func BenchmarkOverlayPerformance(bt *testing.B) {

	o, err := NewLocalOverlay(filepath.Join(bt.TempDir(), "overlay_root"), NewJsonMarshal[OverlayMetadata]())
	if err != nil {
		bt.Errorf("Cannot create test folder: %s", err)
		return
	}

	fileSize := 1024 // 12 MB
	for b := range bt.N {
		out, err := o.OpenWrite(fmt.Sprintf("file %d", b))
		if err != nil {
			bt.Errorf("Opening file for write failed: %s", err)
			return
		}
		if _, err = out.Write(make([]byte, fileSize)); err != nil {
			bt.Errorf("Writing to file failed: %s", err)
			return
		}

		if err = out.Close(); err != nil {
			bt.Errorf("Closing file after write failed: %s", err)
			return
		}
		bt.SetBytes(int64(fileSize))
	}
}
