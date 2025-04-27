package directory

import (
	"strconv"
	"testing"

	"github.com/arsnazarenko/go-hashdb/hashdb/page"
)

func TestDirectorySplit(t *testing.T) {
	t.Run("Split test", func(t *testing.T) {
		p := page.PageFrom(make([]byte, page.PAGE_SIZE))
		hashes := []uint{}
		keys := [][]byte{}

		for i := 1000; ; i++ {
			k := []byte(strconv.Itoa(i))
			if err := p.Put(k, k); err != nil {
				break
			}
			h := hash(k)
			hashes = append(hashes, h)
			keys = append(keys, k)
		}
		d := Directory{
			Meta: Meta{
				Table:      []int{},
				Gd:         1,
				LastPageId: 0,
			},
			DM: nil,
		}
		_, _ = d.split(p)
	})
}
