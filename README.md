# go-hashdb

HashDB -- simple embedded key-value store based on extendible hashing algorithm, implemented in Golang

### Roadmap
* [ ] Add debug print functions (with Stringer) for hashdb 
* [ ] Errors 
* [ ] Update metdata file after each completed `Put` operation for durabilty (current impl update metadata file only after `Close()`)
* [ ] Impl ThreadSafe hashdb impl with `RWLock`
* [ ] Add generic hashdb with support for custom key and value types

