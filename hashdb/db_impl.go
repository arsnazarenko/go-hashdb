package hashdb

var _ DB = (*HashDb)(nil)

type HashDb struct {
	// TODO:
}

// Close implements DB.
func (h *HashDb) Close() error {
	panic("unimplemented")
}

// Get implements DB.
func (h *HashDb) Get(key string) (string, error) {
	panic("unimplemented")
}

// Put implements DB.
func (h *HashDb) Put(key string, value string) error {
	panic("unimplemented")
}

func Open(path string) *HashDb {
    //  TODO:
    return nil
}

