package record

import "encoding/binary"

type ByteRecord []byte

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
	RECORD_KEY_LEN_SZ      = 2                                     // Size of the value len in record header
	RECORD_TOTAL_HEADER_SZ = RECORD_KEY_LEN_SZ + RECORD_VAL_LEN_SZ // Total size of the record header
)

func (r ByteRecord) Key() []byte {
	l := len(r)
	key_start_idx := l - RECORD_TOTAL_HEADER_SZ - int(r.ValueLen()) - int(r.KeyLen())
	key_end_idx := l - RECORD_TOTAL_HEADER_SZ - int(r.ValueLen())
	return r[key_start_idx:key_end_idx]
}

func (r ByteRecord) Value() []byte {
	l := len(r)
	return r[l-RECORD_TOTAL_HEADER_SZ-int(r.ValueLen()) : l-RECORD_TOTAL_HEADER_SZ]
}

func (r ByteRecord) KeyLen() uint16 {
	l := len(r)
	return binary.LittleEndian.Uint16(r[l-RECORD_TOTAL_HEADER_SZ:])
}

func (r ByteRecord) ValueLen() uint16 {
	l := len(r)
	return binary.LittleEndian.Uint16(r[l-RECORD_VAL_LEN_SZ:])
}

func (r ByteRecord) Write(key, value string) {
	keyLen := uint16(len(key))
	valueLen := uint16(len(value))
	total := valueLen + keyLen
	copy(r[:], key)
	copy(r[keyLen:], value)
	binary.LittleEndian.PutUint16(r[total:], keyLen)
	binary.LittleEndian.PutUint16(r[total+RECORD_KEY_LEN_SZ:], valueLen)
}
