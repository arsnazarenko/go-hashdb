package directory

import (
	"bytes"
	"errors"
	"fmt"
	"hash/fnv"

	"github.com/arsnazarenko/go-hashdb/hashdb/disk"
	"github.com/arsnazarenko/go-hashdb/hashdb/page"
	"github.com/arsnazarenko/go-hashdb/hashdb/util"
)

var errHashOutOfTheDirTable = errors.New("hash value out of page id's list")
var offsetOutOfTheData = errors.New("page offset out of data file")

type Meta struct {
	Table      []int // page id's. use signed int for special error id -1
	Gd         uint  // global depth
	LastPageId int   // last page id
}
type Directory struct {
	Meta Meta
	DM   disk.DiskManager
}

func hash(key []byte) uint {
	hash := fnv.New32a()
	hash.Write(key)
	return uint(hash.Sum32())
}

func (d *Directory) extHash(key []byte) uint {
	return hash(key) & ((1 << d.Meta.Gd) - 1)
}

// return Page with PageId or Error
func (d *Directory) getPage(key []byte) (page.Page, int, error) {
	hash := d.extHash(key)
	if hash > uint(len(d.Meta.Table)-1) { // len of Table always >= 1
		return nil, -1, fmt.Errorf("directory.getPage: %w", errHashOutOfTheDirTable)
	}
	id := d.Meta.Table[hash]
	offset := id * page.PAGE_SIZE

	if offset+page.PAGE_SIZE > len(d.DM.Memory()) {
		return nil, -1, fmt.Errorf("directory.GetPage: %w", offsetOutOfTheData)
	}
	return page.PageFrom(d.DM.Memory()[offset : offset+page.PAGE_SIZE]), id, nil
}

func (d *Directory) Get(key []byte) ([]byte, error) {
	p, _, err := d.getPage(key)
	if err != nil {
		return nil, err
	}

	return p.Get(key)
}

func (d *Directory) expand() {
	d.Meta.Table = append(d.Meta.Table, d.Meta.Table...)
	d.Meta.Gd++
}

func (d *Directory) split(p page.Page) (page.Page, page.Page) {
	util.Assert(p.Ld() < d.Meta.Gd, "Local depth of splited page should be less than director")
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
	d.Meta.LastPageId++
	return d.Meta.LastPageId
}

func (d *Directory) replace(splitedPageId int, ld uint) (int, int) {
	newPageId := d.nextPageId()
	for i := range d.Meta.Table {
		if d.Meta.Table[i] == splitedPageId && ((i>>ld)&0x1) == 1 {
			d.Meta.Table[i] = newPageId
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
			copy(d.DM.Memory()[id*page.PAGE_SIZE:(id*page.PAGE_SIZE)+page.PAGE_SIZE], cleaned)
			return nil
		}
		if !errors.Is(err, page.ErrPageIsFull) {
			return fmt.Errorf("directory.Put: %w", err)
		}
		oldMeta := d.Meta
		if p.Ld() == d.Meta.Gd {
			d.expand()
		}
		if p.Ld() < d.Meta.Gd {
			splPage, newPage := d.split(cleaned)
			splId, newId := d.replace(id, p.Ld())
			splPage.SetLd(uint16(p.Ld()) + 1)
			newPage.SetLd(uint16(p.Ld()) + 1)
			if newId*page.PAGE_SIZE >= len(d.DM.Memory()) {
				if err := d.DM.IncreaseSize(); err != nil {
					// failed to increase => restore old meta
					d.Meta = oldMeta
					return fmt.Errorf("directory.Put: %w", err)
				}
			}
			copy(d.DM.Memory()[splId*page.PAGE_SIZE:splId*page.PAGE_SIZE+page.PAGE_SIZE], splPage)
			copy(d.DM.Memory()[newId*page.PAGE_SIZE:newId*page.PAGE_SIZE+page.PAGE_SIZE], newPage)
			d.put(key, value)
		}

	}

	return nil
}

func (d *Directory) getPageById(id int) page.Page {
	util.Assert(id <= d.Meta.LastPageId, "Invalid page id")
	return page.PageFrom(d.DM.Memory()[id*page.PAGE_SIZE : id*page.PAGE_SIZE+page.PAGE_SIZE])

}
func (d *Directory) Put(key, value []byte) error {
	err := d.put(key, value)
	if err != nil {
		return err
	}
	d.DM.Flush()
	return nil
}

func (d *Directory) String() string {

	b := bytes.NewBuffer(make([]byte, 0))
	b.WriteString(fmt.Sprintf("Gd: %d, Table: %v\n", d.Meta.Gd, d.Meta.Table))
	for i := 0; i <= d.Meta.LastPageId; i++ {
		b.WriteString(fmt.Sprintf("Page %d:\n", i))
		p := d.getPageById(i)
		b.WriteString(p.String() + "\n")
	}
	return b.String()
}
