package storage

import (
	"sync"

	"github.com/dgraph-io/badger/v4"
)

type Database struct {
	name string
	path string

	// badger db
	db *badger.DB

	closed bool
}

func NewDatabase(name, path string) (*Database, error) {
	return &Database{
		name: name,
		path: path,
	}, nil
}

func (d *Database) Open() error {
	if d.closed {
		return nil
	}
	db, err := badger.Open(badger.DefaultOptions(d.path + "/" + d.name))
	if err != nil {
		return err
	}
	d.db = db
	return nil
}

func (d *Database) Close() error {
	if d.closed {
		return nil
	}
	d.closed = true
	return d.db.Close()
}

func (d *Database) Name() string {
	return d.name
}

func (d *Database) Path() string {
	return d.path
}

type Catalog struct {
	db *Database
}

func NewCatalog(db *Database) *Catalog {
	return &Catalog{
		db: db,
	}
}

func (c *Catalog) Init() error {

	return nil
}

type Lock struct {
	mu    sync.Mutex
	locks map[string]bool
}

func NewLocks() *Lock {
	return &Lock{
		locks: make(map[string]bool),
	}
}

func (l *Lock) Lock(key string) {
	l.mu.Lock()
	l.locks[key] = true
	l.mu.Unlock()
}

func (l *Lock) Unlock(key string) {
	l.mu.Lock()
	delete(l.locks, key)
	l.mu.Unlock()
}

func (l *Lock) IsLocked(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.locks[key]
}

type Primary struct {
	coll map[string]*Collection
	lock *Lock
}
