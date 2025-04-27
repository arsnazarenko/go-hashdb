# go-hashdb

HashDB -- simple embedded key-value store based on extendible hashing algorithm, implemented in Golang

## Tests
```
make test
```
## Benchmarks
```
make bench
```

## Run example
```
make run
```
## ✅Roadmap
* [ ] Errors 
* [ ] Update metdata file after each completed `Put` operation for durabilty (current impl update metadata file only after `Close()`)
* [ ] Impl ThreadSafe hashdb impl with `RWLock`
* [ ] Add generic hashdb with support for custom key and value types
* [ ] CI
* [ ] Linters
* [ ] GRPC/REST/CLI interfaces




