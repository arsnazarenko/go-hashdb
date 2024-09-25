package page

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPageIterator(t *testing.T) {
	t.Run("Create iterator from invalid page", func(t *testing.T) {
		p := PageFrom(make([]byte, PAGE_SIZE))
		require.Panics(t, func() { NewPageIterator(p, 4093) })
	})
	t.Run("Iterate over zero page", func(t *testing.T) {
		p := PageFrom(make([]byte, PAGE_SIZE))
		it := NewPageIterator(p, p.Use())
		cnt := 0
		for ; it.HasNext(); it.Next() {
			cnt++
		}
		require.Equal(t, 0, cnt)

	})
	t.Run("Iterate over some page", func(t *testing.T) {
		p := PageFrom(make([]byte, PAGE_SIZE))
		for i := 1000; i < 1300; i++ {
			k, v := []byte(strconv.Itoa(i)), []byte(strconv.Itoa(i))
			p.Put(k, v)
		}
		it := NewPageIterator(p, p.Use())
		expectedKV := 1300
		for i := it; i.HasNext(); i.Next() {
			expectedKV--
			r := i.Get()
			require.Equal(t, []byte(strconv.Itoa(expectedKV)), r.Key())
			require.Equal(t, []byte(strconv.Itoa(expectedKV)), r.Value())
		}
		require.Equal(t, 1000, expectedKV)
	})
}
