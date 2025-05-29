
package storage

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

// Database interface for blockchain storage
type Database interface {
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte) error
	Delete(key []byte) error
	Has(key []byte) (bool, error)
	Close() error
	NewBatch() Batch
}

// Batch interface for batch operations
type Batch interface {
	Put(key []byte, value []byte) error
	Delete(key []byte) error
	Write() error
	Reset()
}

// LevelDB implementation
type LevelDB struct {
	db *leveldb.DB
}

// NewLevelDB creates a new LevelDB instance
func NewLevelDB(path string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open leveldb: %v", err)
	}

	return &LevelDB{db: db}, nil
}

// Get retrieves a value by key
func (ldb *LevelDB) Get(key []byte) ([]byte, error) {
	return ldb.db.Get(key, nil)
}

// Put stores a key-value pair
func (ldb *LevelDB) Put(key []byte, value []byte) error {
	return ldb.db.Put(key, value, nil)
}

// Delete removes a key-value pair
func (ldb *LevelDB) Delete(key []byte) error {
	return ldb.db.Delete(key, nil)
}

// Has checks if a key exists
func (ldb *LevelDB) Has(key []byte) (bool, error) {
	return ldb.db.Has(key, nil)
}

// Close closes the database
func (ldb *LevelDB) Close() error {
	return ldb.db.Close()
}

// NewBatch creates a new batch
func (ldb *LevelDB) NewBatch() Batch {
	return &LevelDBBatch{batch: new(leveldb.Batch), db: ldb.db}
}

// LevelDBBatch implements batch operations for LevelDB
type LevelDBBatch struct {
	batch *leveldb.Batch
	db    *leveldb.DB
}

// Put adds a key-value pair to the batch
func (b *LevelDBBatch) Put(key []byte, value []byte) error {
	b.batch.Put(key, value)
	return nil
}

// Delete adds a delete operation to the batch
func (b *LevelDBBatch) Delete(key []byte) error {
	b.batch.Delete(key)
	return nil
}

// Write commits the batch
func (b *LevelDBBatch) Write() error {
	return b.db.Write(b.batch, nil)
}

// Reset resets the batch
func (b *LevelDBBatch) Reset() {
	b.batch.Reset()
}
