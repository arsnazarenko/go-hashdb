package record

import (
	"encoding/binary"

	"github.com/arsnazarenko/go-hashdb/hashdb/util"
)

/*
	    The record does not own the underlying memory.
		It's just a view over memory
		Record binary representation:

		[key, value, keyLen, valueLen]
		|------------|---------------|
		|   PAYLOAD  | RECORD_HEADER |
*/
const (
	RECORD_VAL_LEN_SZ      = 2                                     // Size of the value len in record header
	RECORD_KEY_LEN_SZ      = 2                                     // Size of the key len in record header
	RECORD_TOTAL_HEADER_SZ = RECORD_KEY_LEN_SZ + RECORD_VAL_LEN_SZ // Total size of the record header
)

type Record []byte

func RecordFrom(mem []byte) Record {
	r := Record(mem)
	util.Assert(
		(int(r.KeyLen()+r.ValueLen())+RECORD_TOTAL_HEADER_SZ) <= len(mem),
		"record.RecordFrom: underlying memory less than required Record memory size")
	return r
}

func (r Record) Key() []byte {
	l := len(r)
	key_start_idx := l - RECORD_TOTAL_HEADER_SZ - int(r.ValueLen()+r.KeyLen())
	key_end_idx := l - RECORD_TOTAL_HEADER_SZ - int(r.ValueLen())
	return r[key_start_idx:key_end_idx]
}

func (r Record) Value() []byte {
	l := len(r)
	return r[l-RECORD_TOTAL_HEADER_SZ-int(r.ValueLen()) : l-RECORD_TOTAL_HEADER_SZ]
}

func (r Record) KeyLen() uint16 {
	l := len(r)
	return binary.LittleEndian.Uint16(r[l-RECORD_TOTAL_HEADER_SZ:])
}

func (r Record) ValueLen() uint16 {
	l := len(r)
	return binary.LittleEndian.Uint16(r[l-RECORD_VAL_LEN_SZ:])
}

func (r Record) Write(key, value []byte) {
	keyLen := uint16(len(key))
	valueLen := uint16(len(value))
	total := valueLen + keyLen
	copy(r[:], key)
	copy(r[keyLen:], value)
	binary.LittleEndian.PutUint16(r[total:], keyLen)
	binary.LittleEndian.PutUint16(r[total+RECORD_KEY_LEN_SZ:], valueLen)
}

func (r Record) String() string {
	return "{ " + string(r.Key()) + ", " + string(r.Value()) + " }"
}
