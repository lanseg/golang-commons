package almostio

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

var (
	nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9\.]+`)
)

const (
	maxNameLength         = 200
	mimeBlockSize         = 512
	defaultDirPermissions = 0777
	defaultPermissions    = 0644
	systemFolderName      = ".overlay"
	metadataFileName      = "metadata.json"
)

type FileMetadata struct {
	Name      string `json:"name"`
	LocalName string `json:"localName"`
	Sha256    string `json:"sha256"`
	Mime      string `json:"mime"`
}

// Overlay is an extra layer between the filesystem (io) and the user code.
type Overlay interface {
	OpenRead(name string) (io.ReadCloser, error)
	OpenWrite(name string) (io.WriteCloser, error)
}

// OverlayMetadata contains a system information for the overlay (e.g. file list)
type OverlayMetadata struct {
	FileMetadata map[string]*FileMetadata `json:"fileMetadata"`
}

type localOverlay struct {
	Overlay

	lock sync.Mutex

	marshal  *Marshaller[OverlayMetadata]
	metadata *OverlayMetadata

	root string
}

func (lo *localOverlay) resolve(path ...string) string {
	return filepath.Join(append([]string{lo.root}, path...)...)
}

func (lo *localOverlay) safeName(name string) string {
	h := fnv.New32a()
	h.Write([]byte(name))
	nameHash := fmt.Sprintf("%x", h.Sum([]byte{}))

	newName := nonAlphanumericRegex.ReplaceAllString(name, "_")
	if len(newName) > maxNameLength-len(nameHash)-1 {
		newName = newName[len(newName)-maxNameLength+len(nameHash):]
	}
	return nameHash + "_" + newName
}

func (lo *localOverlay) saveMetadata(fmd *FileMetadata) error {
	lo.lock.Lock()
	defer lo.lock.Unlock()

	if lo.metadata == nil {
		lo.metadata = &OverlayMetadata{
			FileMetadata: map[string]*FileMetadata{},
		}
	}

	if lo.metadata.FileMetadata == nil {
		lo.metadata.FileMetadata = map[string]*FileMetadata{}
	}

	lo.metadata.FileMetadata[fmd.Name] = fmd
	data, err := lo.marshal.Marshal(lo.metadata)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(lo.root, systemFolderName, metadataFileName), data, defaultPermissions)
}

func (lo *localOverlay) OpenRead(name string) (io.ReadCloser, error) {
	f, err := os.Open(lo.resolve(lo.safeName(name)))
	if err != nil {
		return nil, err
	}
	return io.ReadCloser(f), nil
}

func (lo *localOverlay) OpenWrite(name string) (io.WriteCloser, error) {
	lo.lock.Lock()
	defer lo.lock.Unlock()

	safe := lo.safeName(name)
	fwc, err := os.OpenFile(lo.resolve(safe), os.O_WRONLY|os.O_CREATE, defaultPermissions)
	if err != nil {
		return nil, err
	}

	mimeBuffer := bytes.NewBuffer([]byte{})
	sha := sha256.New()
	return NewMultiWriteCloser().
		AddWriteCloser(fwc).
		AddWriter(sha).
		AddWriter(FixedSizeWriter(mimeBuffer, mimeBlockSize)).
		SetOnClose(func() {
			lo.saveMetadata(&FileMetadata{
				Name:      name,
				LocalName: safe,
				Sha256:    fmt.Sprintf("%x", sha.Sum(nil)),
				Mime:      http.DetectContentType(mimeBuffer.Bytes()),
			})
		}), nil
}

func NewLocalOverlay(root string, marshaller *Marshaller[OverlayMetadata]) (Overlay, error) {
	systemFolder := filepath.Join(root, systemFolderName)
	metadataFile := filepath.Join(systemFolder, metadataFileName)
	if err := os.MkdirAll(systemFolder, defaultDirPermissions); err != nil && err != os.ErrExist {
		return nil, err
	}

	if _, err := os.Stat(metadataFile); os.IsNotExist(err) {
		if err = os.WriteFile(metadataFile, []byte("{}"), defaultPermissions); err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(metadataFile)
	if err != nil {
		return nil, err
	}

	mdata, err := marshaller.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	ol := &localOverlay{
		root:     root,
		lock:     sync.Mutex{},
		marshal:  marshaller,
		metadata: mdata,
	}
	return ol, nil
}
