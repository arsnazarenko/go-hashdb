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
    for i := 0; i < 341; i++ {
        _ = db.Put([]byte(strconv.Itoa(tmp)), []byte(strconv.Itoa(tmp)))
        tmp++
    }
    
    tmp = 1000
    for i := 0; i < 341; i++ {
        k := []byte(strconv.Itoa(tmp))
        v, err := db.Get(k)
        if err != nil || string(k) != string(v) {
            panic(err)
        }
        tmp++
    }

}
