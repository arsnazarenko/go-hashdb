package disk

var _ DiskManager = (*ramDiskManager)(nil)

type ramDiskManager struct{}

// Close implements DiskManager.
func (r *ramDiskManager) Close() error {
	panic("unimplemented")
}

func NewRamDiskManager(initCap uint) DiskManager {
	return &ramDiskManager{}
}

// Flush implements DiskManager.
func (r *ramDiskManager) Flush() error {
	panic("unimplemented")
}

// IncreaseSize implements DiskManager.
func (r *ramDiskManager) IncreaseSize() error {
	panic("unimplemented")
}

// Memory implements DiskManager.
func (r *ramDiskManager) Memory() []byte {
	panic("unimplemented")
}
