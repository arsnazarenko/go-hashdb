package hashdb

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/arsnazarenko/go-hashdb/hashdb/directory"
	"github.com/arsnazarenko/go-hashdb/hashdb/disk"
	"github.com/arsnazarenko/go-hashdb/hashdb/page"
)

var _ DB = (*hashDb)(nil)

type hashDb struct {
	dir  *directory.Directory
	path string
}

const (
	FILE_NAME      = "hashdb.data"
	META_FILE_NAME = "hashdb.meta"
)

func Open(path string) (DB, error) {
	if path == "" {
		return nil, fmt.Errorf("hashdb.Open: db direcotry path can't be void")
	}
	var meta directory.Meta

	err := os.MkdirAll(path, 0777)
	if err != nil {
		return nil, fmt.Errorf("hashdb.Open: failed to open datafile dir: %w", err)
	}

	dataFileName := filepath.Join(path, FILE_NAME)
	f, err := os.OpenFile(dataFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("hashdb.Open: failed to open datafile: %w", err)
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("hashdb.Open: %w", err)
	}
	if stat.Size() == 0 { // File not exists, new was created
		meta = directory.Meta{
			Table:      []int{0}, // only one page with PageId = 0
			Gd:         0,
			LastPageId: 0,
		}
		f.Write(make([]byte, page.PAGE_SIZE))
	} else {
		metaRow, err := os.ReadFile(filepath.Join(path, META_FILE_NAME))
		if err != nil {
			return nil, fmt.Errorf("hashdb.Open: read meta file of existing datafile failed: %w", err)
		}
		meta = decodeMeta(bytes.NewBuffer(metaRow))
	}
	dm, err := disk.NewMmapDiskManager(f)
	if err != nil {
		return nil, fmt.Errorf("hashdb.Open: %w", err)
	}
	store := hashDb{
		dir: &directory.Directory{
			Meta: meta,
			DM:   dm,
		},
		path: path,
	}
	return &store, nil
}

// Close implements DB.
func (h *hashDb) Close() error {
	if err := h.dir.DM.Close(); err != nil {
		return fmt.Errorf("hashdb.Close: %w", err)
	}

	metaFile, err := os.OpenFile(filepath.Join(h.path, META_FILE_NAME), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer metaFile.Close()
	if err != nil {
		return fmt.Errorf("hashdb.Close: %w", err)
	}
	metaRaw := encodeMeta(h.dir.Meta).Bytes()
	wr, err := metaFile.WriteAt(metaRaw, 0)
	if err != nil || wr != len(metaRaw) {
		return fmt.Errorf("hashdb.Close: %w", err)
	}
	return nil
}

// Get implements DB.
func (h *hashDb) Get(key []byte) ([]byte, error) {
	return h.dir.Get(key)
}

// Put implements DB.
func (h *hashDb) Put(key []byte, value []byte) error {
	if len(key)+len(value) > 2045 {
		return errors.New("hashDb.Put: the record is too long")
	}
	return h.dir.Put(key, value)
}

func decodeMeta(b *bytes.Buffer) directory.Meta {
	dec := gob.NewDecoder(b)
	var meta directory.Meta
	dec.Decode(&meta)
	return meta
}

func encodeMeta(meta directory.Meta) *bytes.Buffer {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(&meta)
	return &b
}

func (h *hashDb) String() string {
	return h.dir.String()
}
