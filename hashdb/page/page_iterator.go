package page

import (
	"github.com/arsnazarenko/go-hashdb/hashdb/record"
	"github.com/arsnazarenko/go-hashdb/hashdb/util"
)

type PageIterator struct {
	p       Page
	current uint
}

func NewPageIterator(p Page, current uint) *PageIterator {
	util.Assert(current <= PAGE_LOCAL_DEPTH_OFFSET, "position should be less than max page payload size")

	return &PageIterator{p: p, current: current}
}

func (it *PageIterator) HasNext() bool {
	if record.RECORD_TOTAL_HEADER_SZ <= it.current && it.current <= it.p.Use() {
		return true
	}
	return false
}

func (it *PageIterator) Next() {
	br := record.RecordFrom(it.p[:it.current])
	it.current -= (record.RECORD_TOTAL_HEADER_SZ + uint(br.KeyLen()+br.ValueLen()))
}

func (it *PageIterator) Get() record.Record {
	return record.RecordFrom(it.p[:it.current])
}
