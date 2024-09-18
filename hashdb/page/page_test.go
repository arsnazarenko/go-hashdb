package page

import (
	"encoding/binary"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPageFrom(t *testing.T) {
	t.Run("From zero slice", func(t *testing.T) {
		require.Panics(t, func() { PageFrom(make([]byte, 0)) })
	})
	t.Run("From small slice", func(t *testing.T) {
		require.Panics(t, func() { PageFrom(make([]byte, 4090)) })
	})
	t.Run("From memory with invalid Use", func(t *testing.T) {
		require.Panics(t, func() {
			m := make([]byte, 4090)
			binary.LittleEndian.PutUint16(m[PAGE_USE_OFFSET:], math.MaxUint16)
			PageFrom(m)
		})
	})
}

func TestPage(t *testing.T) {
	tests := []struct {
		name   string
		values []int
	}{
		{"Page without values", []int{}},
		{"Simple values", []int{1001, 1002, 1003, 1004, 1005, 1006}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PageFrom(make([]byte, PAGE_SIZE))
			require.Equal(t, uint(0), p.Use())
			require.Equal(t, uint(0), p.ld())
			require.Equal(t, uint(PAGE_LOCAL_DEPTH_OFFSET), p.rest())
			p.setLd(2)
			require.Equal(t, uint(2), p.ld())
			totalLen := 0
			for _, i := range tt.values {
				k, v := strconv.Itoa(i), strconv.Itoa(i)
				totalLen += len(k) + len(v) + 4
				_ = p.Put(k, v)
			}
			require.Equal(t, uint(totalLen), p.Use())
			for _, i := range tt.values {
				actual, _ := p.Get(strconv.Itoa(i))
				require.Equal(t, strconv.Itoa(i), actual)
			}
		})
	}
}
