package record

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var memory = [...]byte{'H', 'E', 'L', 'L', 'O', 'W', 'O', 'R', 'L', 'D', 0x5, 0x0, 0x5, 0x0}
var memory_with_prefix = [...]byte{0xD, 0xE, 0xA, 0xD, 'H', 'E', 'L', 'L', 'O', 'W', 'O', 'R', 'L', 'D', 0x5, 0x0, 0x5, 0x0}
var zero_mem = [...]byte{'H', 'E', 'L', 'L', 'O', 'W', 'O', 'R', 'L', 'D', 0x0, 0x0, 0x0, 0x0}
var onlyHeaderMem = [...]byte{0x0, 0x0, 0x0, 0x0}

func Test(t *testing.T) {
	mem := memory[:]
	br := RecordFrom(mem)
	require.Equal(t, []byte("HELLO"), br.Key())
	require.Equal(t, []byte("WORLD"), br.Value())
	require.Equal(t, uint16(5), br.KeyLen())
	require.Equal(t, uint16(5), br.ValueLen())
}

func TestByteRecordFrom(t *testing.T) {
   t.Run("Create record from mem with only header", func(t *testing.T) {
       require.Panics(t, func(){RecordFrom(onlyHeaderMem[:])})
   }) 
   t.Run("Create record from zero len slice", func(t *testing.T) {
       require.Panics(t, func(){RecordFrom([]byte{})})
   }) 
}

func TestByteRecord(t *testing.T) {
	tests := []struct {
		name           string
		mem            []byte
		expected_key   string
		expected_value string
	}{
		{name: "Key and value", mem: memory[:], expected_key: "HELLO", expected_value: "WORLD"},
		{name: "Key and value with prefix", mem: memory_with_prefix[:], expected_key: "HELLO", expected_value: "WORLD"},
		{name: "Key and value with zero len", mem: zero_mem[:], expected_key: "", expected_value: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := RecordFrom(tt.mem)
			require.Equal(t, []byte(tt.expected_key), br.Key())
			require.Equal(t, []byte(tt.expected_value), br.Value())
			require.Equal(t, uint16(len(tt.expected_key)), br.KeyLen())
			require.Equal(t, uint16(len(tt.expected_value)), br.ValueLen())
		})
	}
}
