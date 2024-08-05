package hashdb


type HashDB interface {
	Close() error
	Put(key, value string) error
	Get(key string) (string, error)
	// TODO: Iterate()
	// TODO: DebugPrint()
}

