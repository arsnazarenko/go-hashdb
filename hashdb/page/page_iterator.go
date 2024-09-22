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

// TODO: change to only Nest method, if can't seek => return nil, else => Page
func (it *PageIterator) HasNext() bool {
	if record.RECORD_TOTAL_HEADER_SZ < it.current && it.current <= it.p.Use() {
		return true
	}
	return false
}

func (it *PageIterator) Next() record.Record {
	br := record.RecordFrom(it.p[:it.current])
	util.Assert(
		it.current >= (uint(record.RECORD_TOTAL_HEADER_SZ)+uint(br.KeyLen()+br.ValueLen())),
		"invalid iterator state: current iterator position less then record size")

	it.current -= (uint(record.RECORD_TOTAL_HEADER_SZ) + uint(br.KeyLen()+br.ValueLen()))
	return br
}
