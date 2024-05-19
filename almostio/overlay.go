package almostio

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

var (
	nonAlphanumericRegex = regexp.MustCompile("[^a-zA-Z0-9\\.]+")
)

const (
	maxNameLength         = 250
	mimeBlockSize         = 512
	defaultDirPermissions = 0777
	defaultPermissions    = 0644
	systemFolderName      = ".overlay"
	metadataFileName      = "metadata.json"
)

type FileMetadata struct {
	BucketId  string `json:"bucketId"`
	Name      string `json:"name"`
	LocalName string `json:"localName"`
	Sha256    string `json: "sha256"`
	Mime      string `json: "mime"`
}

func safeName(name string) string {
	newName := nonAlphanumericRegex.ReplaceAllString(name, "_")
	if len(newName) > maxNameLength {
		newName = newName[:len(newName)-maxNameLength]
	}
	return newName
}

type Overlay interface {
	OpenRead(bucketId string, name string) (io.ReadCloser, error)
	OpenWrite(bucketId string, name string) (io.WriteCloser, error)
}

type OverlayMetadata struct {
	FileMetadata map[string](map[string]*FileMetadata) `json:"fileMetadata"`
}

type localOverlay struct {
	Overlay

	lockMap  sync.Mutex
	locks    map[string]*sync.Mutex
	marshal  func(i interface{}) ([]byte, error)
	metadata *OverlayMetadata

	root string
}

func (lo *localOverlay) getLock(key string) *sync.Mutex {
	lo.lockMap.Lock()
	defer lo.lockMap.Unlock()

	mtx, ok := lo.locks[key]
	if ok {
		return mtx
	}

	lo.locks[key] = &sync.Mutex{}
	return lo.locks[key]
}

func (lo *localOverlay) lockBucket(id string) {
	lo.getLock(id).Lock()
}

func (lo *localOverlay) unlockBucket(id string) {
	lo.getLock(id).Unlock()
}

func (lo *localOverlay) resolve(path ...string) string {
	return filepath.Join(append([]string{lo.root}, path...)...)
}

func (lo *localOverlay) saveMetadata(md *FileMetadata) error {
	lo.lockMap.Lock()
	defer lo.lockMap.Unlock()

	if lo.metadata.FileMetadata == nil {
		lo.metadata.FileMetadata = map[string](map[string]*FileMetadata){}
	}

	if lo.metadata.FileMetadata[md.BucketId] == nil {
		lo.metadata.FileMetadata[md.BucketId] = map[string]*FileMetadata{}
	}
	lo.metadata.FileMetadata[md.BucketId][md.Name] = md

	data, err := lo.marshal(lo.metadata)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(lo.root, systemFolderName, metadataFileName), data, defaultPermissions)
}

func (lo *localOverlay) OpenRead(bucketId string, name string) (io.ReadCloser, error) {
	f, err := os.Open(lo.resolve(bucketId, safeName(name)))
	if err != nil {
		return nil, err
	}
	return io.ReadCloser(f), nil
}

func (lo *localOverlay) OpenWrite(bucketId string, name string) (io.WriteCloser, error) {
	lo.lockBucket(bucketId)
	defer lo.unlockBucket(bucketId)

	if err := os.MkdirAll(lo.resolve(bucketId), defaultDirPermissions); err != nil {
		return nil, err
	}
	safe := safeName(name)
	fwc, err := os.OpenFile(
		lo.resolve(bucketId, safe), os.O_WRONLY|os.O_CREATE, defaultPermissions)
	if err != nil {
		return nil, err
	}

	metadata := &FileMetadata{
		BucketId:  bucketId,
		Name:      name,
		LocalName: safe,
		Sha256:    "",
		Mime:      "",
	}

	mimeBuffer := bytes.NewBuffer([]byte{})
	sha := sha256.New()
	return NewMultiWriteCloser().
		AddWriteCloser(fwc).
		AddWriter(sha).
		AddWriter(FixedSizeWriter(mimeBuffer, mimeBlockSize)).
		SetOnClose(func() {
			metadata.Sha256 = fmt.Sprintf("%x", sha.Sum(nil))
			metadata.Mime = http.DetectContentType(mimeBuffer.Bytes())
			lo.saveMetadata(metadata)
		}), nil
}

func NewLocalOverlay(root string) (Overlay, error) {
	systemFolder := filepath.Join(root, systemFolderName)
	metadataFile := filepath.Join(systemFolder, metadataFileName)
	if _, err := os.Stat(root); os.IsNotExist(err) {
		if err = os.MkdirAll(systemFolder, defaultDirPermissions); err != nil {
			return nil, err
		}
		if err = os.WriteFile(metadataFile, []byte("{}"), defaultPermissions); err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(metadataFile)
	if err != nil {
		return nil, err
	}

	mdata := &OverlayMetadata{}
	if err = json.Unmarshal(data, mdata); err != nil {
		return nil, err
	}

	ol := &localOverlay{
		root:     root,
		locks:    map[string]*sync.Mutex{},
		marshal:  json.Marshal,
		metadata: mdata,
	}
	return ol, nil
}
