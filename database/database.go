package database

type DB interface {
	Get([]byte) ([]byte, error)
	Put([]byte, []byte) error
	Delete([]byte) error
	Close() error
	Iterator() Iterator
}

type Iterator interface {
	Next() bool
	Key() []byte
	Value() []byte
	Release()
	Error() error
	Seek([]byte) bool
}
