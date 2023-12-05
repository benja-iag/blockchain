package database

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB struct {
	db *leveldb.DB
}

type LevelDBIterator interface {
	First() bool
	Next() bool
	Last() bool
	Seek(key []byte) bool
	Prev() bool
}

func NewLevelDB(path string) (DB, error) {
	leveldb, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDB{db: leveldb}, nil
}

func (db *LevelDB) Get(key []byte) ([]byte, error) {
	return db.db.Get(key, nil)
}

func (db *LevelDB) Put(key []byte, value []byte) error {
	return db.db.Put(key, value, nil)
}

func (db *LevelDB) Delete(key []byte) error {
	return db.db.Delete(key, nil)
}

func (db *LevelDB) Close() error {
	return db.db.Close()
}

func (db *LevelDB) Iterator() Iterator {
	return db.db.NewIterator(nil, nil)
}
