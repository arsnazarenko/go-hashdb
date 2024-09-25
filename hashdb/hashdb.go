package hashdb

var _ DB = (*hashDb)(nil)

type hashDb struct {
	// TODO:
}

func Open(path string) (DB, error) {
	return nil, nil
}

// Close implements DB.
func (h *hashDb) Close() error {
	panic("unimplemented")
}

// Get implements DB.
func (h *hashDb) Get(key []byte) ([]byte, error) {
	panic("unimplemented")
}

// Put implements DB.
func (h *hashDb) Put(key []byte, value []byte) error {
	panic("unimplemented")
}
