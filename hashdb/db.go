package hashdb

type DB interface {
	Close()
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
	// TODO: Iterate()
	// TODO: DebugPrint()
}
