package page

import (
	"bytes"
	"encoding/binary"
	"errors"

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

var errKeyNotFound error = errors.New("page.Get: key not found")
var errPageIsFull error = errors.New("page.Put: page is full")

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

func (p Page) rest() uint {
	return PAGE_LOCAL_DEPTH_OFFSET - p.Use()
}

func (p Page) ld() uint {
	return uint(binary.LittleEndian.Uint16(p[PAGE_LOCAL_DEPTH_OFFSET:]))
}

func (p Page) setLd(ld uint16) {
	binary.LittleEndian.PutUint16(p[PAGE_LOCAL_DEPTH_OFFSET:], ld)
}

func (p Page) Get(key string) (string, error) {
	it := PageIterator{
		p:       p,
		current: p.Use(),
	}
	keyLen := uint16(len(key))
	for it.HasNext() {
		r := it.Next()
		if keyLen == r.KeyLen() {
			if bytes.Compare(r.Key(), []byte(key)) == 0 {
				return string(r.Value()), nil
			}
		}
	}
	return "", errKeyNotFound
}

func (p Page) Put(key, value string) error {
    // TODO: add checking for the same {key, value}. In this case we can ommit put operation
	payload := uint(len(key) + len(value) + record.RECORD_TOTAL_HEADER_SZ)
	if p.rest() >= payload {
		use := p.Use()
		r := record.RecordFrom(p[use : use+payload])
		r.Write(key, value)
		binary.LittleEndian.PutUint16(p[PAGE_USE_OFFSET:], uint16(use+payload))
		return nil
	}
	return errPageIsFull
}

func (p Page) Gc() Page {
	tmp := PageFrom(make([]byte, PAGE_SIZE))
	tmp.setLd(uint16(p.ld()))
	lookup := make(map[string]bool)
	it := &PageIterator{
		p:       p,
		current: p.Use(),
	}
	for it.HasNext() {
		r := it.Next()
		k := string(r.Key())
		if _, ok := lookup[k]; !ok {
			lookup[k] = true
			tmp.Put(k, string(r.Value()))
		} else {
			continue
		}

	}
	return tmp
}
