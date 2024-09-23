package page

import (
	"encoding/binary"
	"errors"
	"math"
	"strconv"
	"testing"

	"github.com/arsnazarenko/go-hashdb/hashdb/record"
	"github.com/stretchr/testify/require"
)

func TestPageFrom(t *testing.T) {
	t.Run("From zero slice", func(t *testing.T) {
		require.Panics(t, func() { PageFrom(make([]byte, 0)) })
	})
	t.Run("From small slice", func(t *testing.T) {
		require.Panics(t, func() { PageFrom(make([]byte, 4090)) })
	})
	t.Run("From memory with invalid Use value", func(t *testing.T) {
		require.Panics(t, func() {
			m := make([]byte, 4096)
			binary.LittleEndian.PutUint16(m[PAGE_USE_OFFSET:], math.MaxUint16)
			PageFrom(m)
		})
	})
}

func TestPagePut(t *testing.T) {
    t.Run("Put in empty Page", func (t *testing.T) {
        page := PageFrom(make([]byte, 4096))
        require.Equal(t, page.Use(), uint(0))
        k, v := []byte("Hello"), []byte("world")
        newLen := len(k) + len(v) + record.RECORD_TOTAL_HEADER_SZ
        err := page.Put(k, v)
        require.Nil(t, err)
        require.Equal(t, page.Use(), uint(newLen))
        actualV, err := page.Get(k)
        require.Nil(t, err)
        require.Equal(t, v, actualV)
    })
    t.Run("Put in full Page", func (t *testing.T) {
		m := make([]byte, 4096)
		binary.LittleEndian.PutUint16(m[PAGE_USE_OFFSET:], PAGE_LOCAL_DEPTH_OFFSET - 1)
        p := PageFrom(m)
        k, v := []byte("Hello"), []byte("world")
        err := p.Put(k, v)
        require.ErrorIs(t, errPageIsFull, errors.Unwrap(err))

    })
    t.Run("Put value with same key full Page", func (t *testing.T) {
		m := make([]byte, 4096)
        p := PageFrom(m)
        k1, v1, v2 := []byte("Hello"), []byte("world1"), []byte("world2")
        _ = p.Put(k1, v1)
        actualV1, _ := p.Get(k1)
        require.Equal(t, v1, actualV1)
        _ = p.Put(k1, v2)
        actualV2, _ := p.Get(k1)
        require.Equal(t, v2, actualV2)
    })
}
func TestPageGet(t *testing.T) {
    
    t.Run("Get from empty Page", func (t *testing.T) {
		m := make([]byte, 4096)
        p := PageFrom(m)
        k := []byte("Hello")
        _, err := p.Get(k)
        require.ErrorIs(t, errKeyNotFound, errors.Unwrap(err))
    })
}
func TestPageGc(t *testing.T) {
		m := make([]byte, 4096)
        p := PageFrom(m)
        k := []byte("Hello")
        vals :=[][]byte { []byte("Hello"), []byte("world1"), []byte("world2"), []byte("world3")}
        for _, v := range vals {
            p.Put(k, v)
        }
        p.Gc()
        require.Equal(t, p.Use(), uint(len(k) + len(vals[len(vals) - 1]) + record.RECORD_TOTAL_HEADER_SZ))
        v, err := p.Get(k)
        require.Nil(t, err)
        require.Equal(t, v, vals[len(vals) - 1])
}


func TestPutScenarios(t *testing.T) {
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
			require.Equal(t, uint(0), p.Ld())
			require.Equal(t, uint(PAGE_LOCAL_DEPTH_OFFSET), p.rest())
			p.setLd(2)
            require.Equal(t, uint(2), p.Ld())
			totalLen := 0
			for _, i := range tt.values {
				k, v := strconv.Itoa(i), strconv.Itoa(i)
				totalLen += len(k) + len(v) + 4
				_ = p.Put([]byte(k), []byte(v))
			}
			require.Equal(t, uint(totalLen), p.Use())
			for _, i := range tt.values {
				actual, _ := p.Get([]byte(strconv.Itoa(i)))
				require.Equal(t, []byte(strconv.Itoa(i)), actual)
			}
		})
	}
}
