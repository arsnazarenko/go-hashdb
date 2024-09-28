package bench

import (
	"crypto/rand"
	"fmt"
	"os"

	"github.com/arsnazarenko/go-hashdb/hashdb"
)

const dirparh = "/tmp/hashdb_bench"

var db hashdb.DB = nil

func generateValue(valueLen int) []byte {
	v := make([]byte, valueLen)
	_, err := rand.Read(v)
	if err != nil {
		panic(err)
	}
	return v
}

func reset() {
	if db != nil {
		_ = db.Close()
	}
	err := os.RemoveAll(dirparh)
	if err != nil {
		panic(err)
	}
}

func open() {
	tmpDb, err := hashdb.Open(dirparh)
	if err != nil {
		panic(err)
	}
	db = tmpDb
}

func openAndFill(count int, value []byte) {
	open()
	for i := 0; i < count; i++ {
		if err := db.Put([]byte(fmt.Sprintf("key_%d", i)), value); err != nil {
			panic(err)
		}
	}
}
