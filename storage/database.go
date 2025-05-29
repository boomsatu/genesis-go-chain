
package storage

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// Database interface for blockchain storage
type Database interface {
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte) error
	Delete(key []byte) error
	Has(key []byte) (bool, error)
	Close() error
	NewBatch() Batch
	Stats() map[string]string
}

// Batch interface for batch operations
type Batch interface {
	Put(key []byte, value []byte) error
	Delete(key []byte) error
	Write() error
	Reset()
	Size() int
}

// LevelDBOptions holds configuration for LevelDB
type LevelDBOptions struct {
	CacheSize    int // Cache size in MB
	MaxOpenFiles int // Maximum number of open files
	WriteBuffer  int // Write buffer size in MB
}

// LevelDB implementation
type LevelDB struct {
	db *leveldb.DB
}

// NewLevelDB creates a new LevelDB instance with options
func NewLevelDB(path string, options *LevelDBOptions) (*LevelDB, error) {
	opts := &opt.Options{
		BlockCacheCapacity:     options.CacheSize * 1024 * 1024,  // Convert MB to bytes
		OpenFilesCacheCapacity: options.MaxOpenFiles,
		WriteBuffer:           options.WriteBuffer * 1024 * 1024, // Convert MB to bytes
		CompactionTableSize:   4 * 1024 * 1024, // 4MB
		CompactionTotalSize:   16 * 1024 * 1024, // 16MB
	}

	db, err := leveldb.OpenFile(path, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open leveldb at %s: %v", path, err)
	}

	return &LevelDB{db: db}, nil
}

// Get retrieves a value by key
func (ldb *LevelDB) Get(key []byte) ([]byte, error) {
	data, err := ldb.db.Get(key, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, ErrKeyNotFound
		}
		return nil, fmt.Errorf("leveldb get error: %v", err)
	}
	return data, nil
}

// Put stores a key-value pair
func (ldb *LevelDB) Put(key []byte, value []byte) error {
	if err := ldb.db.Put(key, value, nil); err != nil {
		return fmt.Errorf("leveldb put error: %v", err)
	}
	return nil
}

// Delete removes a key-value pair
func (ldb *LevelDB) Delete(key []byte) error {
	if err := ldb.db.Delete(key, nil); err != nil {
		return fmt.Errorf("leveldb delete error: %v", err)
	}
	return nil
}

// Has checks if a key exists
func (ldb *LevelDB) Has(key []byte) (bool, error) {
	exists, err := ldb.db.Has(key, nil)
	if err != nil {
		return false, fmt.Errorf("leveldb has error: %v", err)
	}
	return exists, nil
}

// Close closes the database
func (ldb *LevelDB) Close() error {
	if err := ldb.db.Close(); err != nil {
		return fmt.Errorf("leveldb close error: %v", err)
	}
	return nil
}

// NewBatch creates a new batch
func (ldb *LevelDB) NewBatch() Batch {
	return &LevelDBBatch{
		batch: new(leveldb.Batch),
		db:    ldb.db,
	}
}

// Stats returns database statistics
func (ldb *LevelDB) Stats() map[string]string {
	stats := make(map[string]string)
	
	if stat, err := ldb.db.GetProperty("leveldb.stats"); err == nil {
		stats["general"] = stat
	}
	
	if stat, err := ldb.db.GetProperty("leveldb.sstables"); err == nil {
		stats["sstables"] = stat
	}
	
	if stat, err := ldb.db.GetProperty("leveldb.blockpool"); err == nil {
		stats["blockpool"] = stat
	}
	
	return stats
}

// LevelDBBatch implements batch operations for LevelDB
type LevelDBBatch struct {
	batch *leveldb.Batch
	db    *leveldb.DB
	size  int
}

// Put adds a key-value pair to the batch
func (b *LevelDBBatch) Put(key []byte, value []byte) error {
	b.batch.Put(key, value)
	b.size++
	return nil
}

// Delete adds a delete operation to the batch
func (b *LevelDBBatch) Delete(key []byte) error {
	b.batch.Delete(key)
	b.size++
	return nil
}

// Write commits the batch
func (b *LevelDBBatch) Write() error {
	if err := b.db.Write(b.batch, nil); err != nil {
		return fmt.Errorf("leveldb batch write error: %v", err)
	}
	return nil
}

// Reset resets the batch
func (b *LevelDBBatch) Reset() {
	b.batch.Reset()
	b.size = 0
}

// Size returns the number of operations in the batch
func (b *LevelDBBatch) Size() int {
	return b.size
}

// Custom errors
var (
	ErrKeyNotFound = fmt.Errorf("key not found")
)
