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
	defer db.Close()
	tmp := 1000

	for range 3410 {
		_ = db.Put([]byte(strconv.Itoa(tmp)), []byte(""))
		tmp++
	}

	tmp = 1000
	for range 3410 {
		k := []byte(strconv.Itoa(tmp))
		v, err := db.Get(k)
		if err != nil || string(v) != "" {
			panic(err)
		}
		tmp++
	}

}
