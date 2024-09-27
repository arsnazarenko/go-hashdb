package disk

type DiskManager interface {
	IncreaseSize() error
	Memory() []byte
	Flush() error
	Close() error
}

