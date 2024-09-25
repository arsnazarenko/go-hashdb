package directory

import (
	"errors"
	"fmt"
	"hash/fnv"
	"os"
	"syscall"

	"github.com/arsnazarenko/go-hashdb/hashdb/page"
	"github.com/arsnazarenko/go-hashdb/hashdb/util"
	"github.com/edsrzf/mmap-go"
)

var errHashOutOfTheDirTable = errors.New("hash value out of page id's list")
var offsetOutOfTheData = errors.New("page offset out of data file")

type Directory struct {
	Table      []int     // page id's. use signed int for special error id -1
	data       mmap.MMap // memory mapped buffer buffer
	gd         uint      // global depth
	dataFile   *os.File  // mmap file
	LastPageId int       // last page id
}

func hash(key []byte) uint {
	hash := fnv.New32a()
	hash.Write(key)
	return uint(hash.Sum32())
}

func (d *Directory) extHash(key []byte) uint {
	return hash(key) & ((1 << d.gd) - 1)
}

// return Page with PageId or Error
func (d *Directory) getPage(key []byte) (page.Page, int, error) {
	hash := d.extHash(key)
	if hash > uint(len(d.Table)-1) { // len of Table always >= 1
		return nil, -1, fmt.Errorf("directory.getPage: %w", errHashOutOfTheDirTable)
	}
	id := d.Table[hash]
	offset := id * page.PAGE_SIZE

	if offset+page.PAGE_SIZE > len(d.data) {
		return nil, -1, fmt.Errorf("directory.GetPage: %w", offsetOutOfTheData)
	}
	return page.PageFrom(d.data[offset : offset+page.PAGE_SIZE]), id, nil
}

func (d *Directory) Get(key []byte) ([]byte, error) {
	p, _, err := d.getPage(key)
	return p, err
}

func (d *Directory) expand() {
	d.Table = append(d.Table, d.Table...)
	d.gd++
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
	util.Assert(p.Ld() < d.gd, "Local depth of splited page should be less than director")
	p1 := page.PageFrom(make([]byte, page.PAGE_SIZE))
	p2 := page.PageFrom(make([]byte, page.PAGE_SIZE))
	for i := page.NewPageIterator(p, p.Use()); i.HasNext(); i.Next() {
		k := i.Get().Key()
		h := d.extHash(k)
		if (h>>p.Ld())&0x1 == 1 {
			p2.Put(k, i.Get().Value())
		} else {
			p1.Put(k, i.Get().Value())
		}
	}
	return p1, p2
}

func (d *Directory) nextPageId() int {
	d.LastPageId++
	return d.LastPageId
}

func (d *Directory) replace(splitedPageId int, ld uint) (int, int) {
	newPageId := d.nextPageId()
	for i := 0; i < len(d.Table); i++ {
		if d.Table[i] == splitedPageId && ((i>>ld)&0x1) == 1 {
			d.Table[i] = newPageId
		}
	}
	return splitedPageId, newPageId

}

func (d *Directory) put(key, value []byte) error {
	p, id, err := d.getPage(key)
	if err != nil {
		return fmt.Errorf("directory.Put: %w", err)
	}
	err = p.Put(key, value)
	if err != nil {
		if !errors.Is(err, page.ErrPageIsFull) {
			return fmt.Errorf("directory.Put: %w", err)
		}
		cleaned := p.Gc()
		err := cleaned.Put(key, value)
		if err == nil {
			copy(d.data[id*page.PAGE_SIZE:(id*page.PAGE_SIZE)+page.PAGE_SIZE], cleaned)
			return nil
		}
		if !errors.Is(err, page.ErrPageIsFull) {
			return fmt.Errorf("directory.Put: %w", err)
		}
		dirCopy := *d
		if p.Ld() == dirCopy.gd {
			dirCopy.expand()
		}
		if p.Ld() < dirCopy.gd {
			splPage, newPage := dirCopy.split(cleaned)
			splId, newId := dirCopy.replace(id, p.Ld())
			splPage.SetLd(uint16(p.Ld()) + 1)
			newPage.SetLd(uint16(p.Ld()) + 1)
			if newId*page.PAGE_SIZE >= len(dirCopy.data) {
				if err := dirCopy.increaseSize(); err != nil {
					return fmt.Errorf("directory.Put: %w", err)
				}
				*d = dirCopy
			}
			copy(d.data[splId*page.PAGE_SIZE:splId*page.PAGE_SIZE+page.PAGE_SIZE], splPage)
			copy(d.data[newId*page.PAGE_SIZE:newId*page.PAGE_SIZE+page.PAGE_SIZE], newPage)
			d.put(key, value)
		}

	}

	return nil
}

func (d *Directory) Put(key, value []byte) error {
	err := d.put(key, value)
	if err != nil {
		return err
	}
	d.data.Flush()
	return nil
}

func (d *Directory) mmapDataFile() error {
	d.data.Unmap()
	data, err := mmap.Map(d.dataFile, mmap.RDWR|syscall.MAP_POPULATE, 0644)
	if err != nil {
		return fmt.Errorf("directory.mmapDataFile: %w", err)
	}
	d.data = data
	return nil
}
