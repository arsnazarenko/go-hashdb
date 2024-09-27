package disk

import (
	"fmt"
	"os"
	"syscall"

	"github.com/edsrzf/mmap-go"
)

var _ DiskManager = (*mmapDiskManager)(nil)

type mmapDiskManager struct {
	dataFile *os.File
	data     mmap.MMap
}

func NewMmapDiskManager(f *os.File) (DiskManager, error) {
    var m mmapDiskManager
    m.dataFile = f
    if err := m.mmapDataFile(); err != nil {
		return nil, fmt.Errorf("mmapDiskManager.New: %w", err)
    }
    return &m, nil
}

// Close implements DiskManager.
func (m *mmapDiskManager) Close() error {
	m.Flush()
    m.data.Unmap()
    m.dataFile.Close()
	return nil
}

// Flush implements DiskManager.
func (m *mmapDiskManager) Flush() error {
	return m.data.Flush()
}

func (m *mmapDiskManager) mmapDataFile() error {
	m.data.Unmap()
	data, err := mmap.Map(m.dataFile, mmap.RDWR|syscall.MAP_POPULATE, 0644)
	if err != nil {
		return fmt.Errorf("mmapDiskManager.mmapDataFile: %w", err)
	}
	m.data = data
	return nil
}

// IncreaseSize implements DiskManager.
func (m *mmapDiskManager) IncreaseSize() error {

	stats, err := m.dataFile.Stat()
	if err != nil {
		return fmt.Errorf("mmapDiskManager.IncreaseSize: %w", err)
	}
	size := stats.Size()
	n, err := m.dataFile.WriteAt(make([]byte, size), int64(size)) // x2 increase of dataFile
	if int64(n) < size || err != nil {
		return fmt.Errorf("mmapDiskManager.IncreaseSize: %w", err)
	}
	// remap increased file in RAM
	if err := m.mmapDataFile(); err != nil {
		return fmt.Errorf("mmapDiskManager.IncreaseSize: %w", err)
	}
	return nil
}

// Memory implements DiskManager.
func (m *mmapDiskManager) Memory() []byte {
	return m.data
}
