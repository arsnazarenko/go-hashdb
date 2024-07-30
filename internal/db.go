package hashdb

type HashDB interface {
    Open(path string) error
    Close() error
    Put(key, value string) error
    Get(key string) (string, error)
    // TODO: Iterate()
    // TODO: DebugPrint()
}
