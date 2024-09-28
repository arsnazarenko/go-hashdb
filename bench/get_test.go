package bench

import (
	"fmt"
	"testing"
)

func BenchmarkGetValue(b *testing.B) {
	valueSets := map[string][]byte{
		"empty":   []byte(""),
		"small":  generateValue(12),
		"medium": generateValue(128),
		"large":  generateValue(512),
	}
	for key, value := range valueSets {
		b.Run(key, func(b *testing.B) {
			benchmarkGetValueWithLen(b, value)
		})
	}
}

func benchmarkGetValueWithLen(b *testing.B, value []byte) {
	openAndFill(b.N, value) // open database and fill with values
	defer reset()           // // close db and remove db directory

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := db.Get([]byte(fmt.Sprintf("key_%d", i)))
		if err != nil {
			b.Fatal()
		}
	}
}
