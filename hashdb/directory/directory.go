package directory

import (
	"errors"
	"fmt"
	"hash/fnv"
	"os"

	"github.com/arsnazarenko/go-hashdb/hashdb/page"
	"github.com/edsrzf/mmap-go"
)

var errHashOutOfTheDirTable = errors.New("hash value out of page id's list")
var offsetOutOfTheData = errors.New("page offset out of data file")

type Directory struct {
	Table      []int     // page id's. use signed int for special error id -1
	data       mmap.MMap // memory mapped buffer buffer
	Gd         uint      // global depth
	dataFile   *os.File  // mmap file
	LastPageId int       // last page id
}

func (d *Directory) extHash(key []byte) uint {
	hash := fnv.New32a()
	hash.Write(key)
	return uint(hash.Sum32()) & ((1 << d.Gd) - 1)
}

// return Page with PageId or Error
func (d *Directory) getPage(key []byte) (page.Page, int, error) {
	hash := d.extHash(key)
	if hash > uint(len(d.Table) - 1) { // len of Table always >= 1
        return nil, -1, fmt.Errorf("directory.getPage: %w", errHashOutOfTheDirTable)
	}
    id := d.Table[hash]
    offset := id * page.PAGE_SIZE

    if (offset + page.PAGE_SIZE > len(d.data)) {
        return nil, -1, fmt.Errorf("directory.GetPage: %w", offsetOutOfTheData)
    }
    return page.PageFrom(d.data[offset:offset + page.PAGE_SIZE]), id, nil
} 

func (d *Directory) Get(key []byte) ([]byte, error) {
    p, _, err := d.getPage(key)
    return p, err
}

func (d *Directory) expand() {
    d.Table = append(d.Table, d.Table...)
    d.Gd++
}

func (d *Directory) increaseSize() error {
    stats, err := d.dataFile.Stat()
    if err != nil {
        return fmt.Errorf("directory.increaseSize: %w", err)
    }
    size := stats.Size()
    n, err := d.dataFile.WriteAt(make([]byte, size), int64(size)) // x2 increase of dataFile
    if int64(n) < size || err != nil {
        return fmt.Errorf("directory.increaseSize: %w", err)
    }
    // remap increased file in RAM
    if err := d.mmapDataFile(); err != nil {
        return fmt.Errorf("directory.increaseSize: %w", err)
    }
    return nil
}


func (d *Directory) split(p page.Page) (page.Page, page.Page) {
    p1 := make([]byte, page.PAGE_SIZE)
    p2 := make([]byte, page.PAGE_SIZE)
    t := page.NewPageIterator(p, p.Use())
    _ = t  
    // IMPLEMENT 
    return p1, p2
}



func (d *Directory) mmapDataFile() error {
    return nil
}


















