package example

import "github.com/arsnazarenko/go-hashdb/hashdb"

func main() {
    db := hashdb.Open("/tmp/hashdb")
    defer db.Close() 
}
