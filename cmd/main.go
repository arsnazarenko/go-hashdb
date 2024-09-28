package main

import (
	"strconv"

	"github.com/arsnazarenko/go-hashdb/hashdb"
)

func main() {
	db, err := hashdb.Open("/tmp/hashdb")
	if err != nil {
		panic(err)
	}
	defer func() { _ = db.Close() }()
	tmp := 1000
	for i := 0; i < 3410; i++ {
		_ = db.Put([]byte(strconv.Itoa(tmp)), []byte(""))
		tmp++
	}

	tmp = 1000
	for i := 0; i < 3410; i++ {
		k := []byte(strconv.Itoa(tmp))
		v, err := db.Get(k)
		if err != nil ||  string(v) != "" {
			panic(err)
		}
		tmp++
	}

}
