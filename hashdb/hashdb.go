package hashdb

var _ DB = (*HashDb)(nil)

type HashDb struct {
	// TODO:
}

func Open(path string) (*HashDb, error) {
    panic("unimplemented")
}

// Close implements DB.
func (h *HashDb) Close() error {
	panic("unimplemented")
}

// Get implements DB.
func (h *HashDb) Get(key []byte) ([]byte, error) {
	panic("unimplemented")
}

// Put implements DB.
func (h *HashDb) Put(key []byte, value []byte) error {
	panic("unimplemented")
}

