package page

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/arsnazarenko/go-hashdb/hashdb/record"
	"github.com/arsnazarenko/go-hashdb/hashdb/util"
)

// Page layout:
//
// [Record1, Record2, Record3, ..., PageN, local_depth(2B), page_use(2B)]

const (
	PAGE_SIZE               = 4096
	PAGE_USE_OFFSET         = 4094
	PAGE_LOCAL_DEPTH_OFFSET = 4092
)

var ErrKeyNotFound error = errors.New("key not found")
var ErrPageIsFull error = errors.New("page is full")

type Page []byte

func PageFrom(mem []byte) Page {
	util.Assert(len(mem) >= PAGE_SIZE, "page.PageFrom: len underlying memory of Page should be greater or equal 4096")
	p := Page(mem)
	util.Assert(p.Use() <= PAGE_LOCAL_DEPTH_OFFSET, "page.PageFrom: max len of page payload is 4092")
	return p
}

func (p Page) Use() uint {
	return uint(binary.LittleEndian.Uint16(p[PAGE_USE_OFFSET:]))
}

func (p Page) Ld() uint {
	return uint(binary.LittleEndian.Uint16(p[PAGE_LOCAL_DEPTH_OFFSET:]))
}

func (p Page) SetLd(ld uint16) {
	binary.LittleEndian.PutUint16(p[PAGE_LOCAL_DEPTH_OFFSET:], ld)
}

func (p Page) SetUse(use uint16) {
	util.Assert(p.Use() <= PAGE_LOCAL_DEPTH_OFFSET, "page.PageFrom: max len of page payload is 4092")
	binary.LittleEndian.PutUint16(p[PAGE_USE_OFFSET:], use)
}

func (p Page) Get(key []byte) ([]byte, error) {
	keyLen := uint16(len(key))
	it := NewPageIterator(p, p.Use())
	for i := it; it.HasNext(); i.Next() {
		r := i.Get()
		if keyLen == r.KeyLen() {
			if bytes.Compare(r.Key(), []byte(key)) == 0 {
				return r.Value(), nil
			}
		}
	}
	return nil, fmt.Errorf("page.Get: %w", ErrKeyNotFound)
}

func (p Page) Put(key, value []byte) error {
	// TODO: add checking for the same {key, value}. In this case we can ommit put operation
	payload := uint(len(key) + len(value) + record.RECORD_TOTAL_HEADER_SZ)
	if p.rest() >= payload {
		use := p.Use()
		r := record.Record(p[use : use+payload]) // not use RecordFrom, because Put does not validate this memory
		r.Write(key, value)
		binary.LittleEndian.PutUint16(p[PAGE_USE_OFFSET:], uint16(use+payload))
		return nil
	}
	return fmt.Errorf("page.Put: %w", ErrPageIsFull)
}

// like compaction, delete older versions of updated k,v
func (p Page) Gc() Page {
	tmp := PageFrom(make([]byte, PAGE_SIZE))
	tmp.SetLd(uint16(p.Ld()))
	lookup := make(map[string]bool)
	for i := NewPageIterator(p, p.Use()); i.HasNext(); i.Next() {
		r := i.Get()
		k := string(r.Key())
		if _, ok := lookup[k]; !ok {
			lookup[k] = true
			tmp.Put(r.Key(), r.Value())
		} else {
			continue
		}
	}
	return tmp
}

func (p Page) rest() uint {
	return PAGE_LOCAL_DEPTH_OFFSET - p.Use()
}

func (p Page) String() string {
	b := bytes.NewBufferString("[")
	cnt := 0
	for it := NewPageIterator(p, p.Use()); it.HasNext(); it.Next() {
		if cnt > 60 {
			b.WriteString("\n")
			cnt = 0
		}
		str := it.Get().String() + " "
		b.WriteString(str)
		cnt += len(str)
	}
	b.WriteString("]\n" + fmt.Sprintf("ld: %d, use: %d", p.Ld(), p.Use()))
	return b.String()
}
