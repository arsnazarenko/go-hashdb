package hashdb

// Constraint for type that can be used for key or value
type Serializable interface {
	Serialize() []byte
	Deserialize(from []byte)
}

type GenericDb[K Serializable, V Serializable] interface {
	Close() error
	Put(key K, value V) error
	Get(key K) (V, error)
}

type GenericHashDb[K Serializable, V Serializable] struct {
	hashdb hashDb
}

func (g *GenericHashDb[K, V]) Close() error {
	return g.hashdb.Close()
}
func (g *GenericHashDb[K, V]) Put(key K, value V) error {
	return g.hashdb.Put(key.Serialize(), value.Serialize())

}
func (g *GenericHashDb[K, V]) Get(key K) (V, error) {
	var v V
	raw, err := g.hashdb.Get(key.Serialize())
	v.Deserialize(raw) // TODO: deserialize from raw check
	return v, err
}
