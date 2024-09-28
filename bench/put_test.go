package bench

import (
	"fmt"
	"testing"
)

func BenchmarkPutValue(b *testing.B) {
	valueSets := map[string][]byte{
		"64B":  generateValue(64),
		"128B": generateValue(128),
		"256B": generateValue(256),
		"512B":  generateValue(512),
	}
	for key, value := range valueSets {
		b.Run(key, func(b *testing.B) {
			benchmarkPutValueWithLen(b, value)
		})
	}
}

func benchmarkPutValueWithLen(b *testing.B, value []byte) {
	open()        // open db
	defer reset() // close db and remove db directory
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := db.Put([]byte(fmt.Sprintf("key_%d", i)), value)
		if err != nil {
			b.Fatal()
		}
	}
}
