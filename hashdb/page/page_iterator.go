package page

import "github.com/arsnazarenko/go-hashdb/hashdb/record"


type PageIterator struct {
    p Page
    current uint
}

func (it *PageIterator) HasNext() bool {
    if record.RECORD_TOTAL_HEADER_SZ < it.current && it.current <= it.p.Use() {
        return true
    }
    return false
}

func (it *PageIterator) Next() record.Record {
    br := record.RecordFrom(it.p[:it.current])
    it.current -= (uint(record.RECORD_TOTAL_HEADER_SZ) + uint(br.KeyLen() + br.ValueLen())) 
    return br
}
