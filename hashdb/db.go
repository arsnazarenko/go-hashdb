package hashdb


type DB interface {
	Close() error
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
	// TODO: Iterate()
	// TODO: DebugPrint()
}

