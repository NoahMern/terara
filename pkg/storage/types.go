package storage

import (
	"errors"

	"github.com/dgraph-io/badger/v4"
	"github.com/noahmern/terara/pkg/types"
)

var (
	ErrStaticDocument = errors.New("static document")
)

// this is not thread safe
type Document struct {
	// this represents the document type
	kv map[string]types.Object

	db   *Database
	coll *Collection
	tnx  *badger.Txn
	key  []byte

	modified bool
	static   bool // if true we can't modify this document
}

var _ types.Document = (*Document)(nil)

func NewDocument(db *Database, coll *Collection, tnx *badger.Txn) *Document {
	// create a new document
	return &Document{
		kv:   make(map[string]types.Object),
		db:   db,
		coll: coll,
		tnx:  tnx,
	}
}

// FIXME: implement
func NewStaticDocument(db *Database) *Document {
	// create a new document
	return &Document{
		kv:     make(map[string]types.Object),
		db:     db,
		static: true,
	}
}

func (d *Document) ID() types.Object {
	// get the id
	return d.kv["id"]
}

func (d *Document) Del(key []byte) error {
	if d.static {
		return ErrStaticDocument
	}
	// delete a key
	delete(d.kv, string(key))
	d.modified = true
	return nil
}

func (d *Document) Get(key []byte) (types.Object, error) {
	// get a value
	return d.kv[string(key)], nil
}

func (d *Document) Set(key []byte, value types.Object) error {
	if d.static {
		return ErrStaticDocument
	}
	// set a value
	d.kv[string(key)] = value
	d.modified = true
	return nil
}

func (d *Document) Keys() [][]byte {
	// get the keys
	keys := make([][]byte, 0)
	for key := range d.kv {
		keys = append(keys, []byte(key))
	}
	return keys
}

func (d *Document) Type() byte {
	// get the type
	return types.DocumentType
}

func (d *Document) Value() interface{} {
	// get the value
	return d.kv
}

func (d *Document) String() string {
	return "Document"
}

func (d *Document) MarshalObject() ([]byte, error) {
	return types.GenericDocumentUnmarshaler(d)
}

func (d *Document) UnmarshalObject(b []byte) (int, error) {
	if len(b) < 1 {
		return 0, types.ErrInvalidLength
	}
	if b[0] != d.Type() {
		return 0, types.ErrInvalidType
	}
	count := 1
	for {
		if len(b) < 1 {
			return 0, types.ErrInvalidLength
		}
		if b[count] == types.EOFType {
			break
		}
		var name types.Name
		nameCount, err := name.UnmarshalObject(b[count:])
		if err != nil {
			return 0, err
		}
		count += nameCount
		value, countValue, err := types.UnmarshalObject(b[count:])
		if err != nil {
			return 0, err
		}
		// value can't be internal values like EOF, Last, Name
		if types.IsInternal(value.Type()) {
			return 0, types.ErrInvalidDocument
		}
		count += countValue
		d.kv[string(name)] = value
	}
	// check for id
	if _, ok := d.kv["id"]; !ok {
		return 0, types.ErrInvalidDocument
	}
	return count, nil
}

func (d *Document) Project(b []byte, keys ...[]byte) (int, error) {
	if len(keys) < 1 {
		return 0, nil
	}
	if len(b) < 1 {
		return 0, types.ErrInvalidLength
	}
	if b[0] != d.Type() {
		return 0, types.ErrInvalidType
	}
	// check if name is in keys
	inKey := func(name types.Name) bool {
		for _, key := range keys {
			if string(name) == string(key) {
				return true
			}
		}
		return false
	}
	count := 1
	for {
		if len(b) < 1 {
			return 0, types.ErrInvalidLength
		}
		if b[count] == types.EOFType {
			break
		}
		var name types.Name
		nameCount, err := name.UnmarshalObject(b[count:])
		if err != nil {
			return 0, err
		}
		count += nameCount
		if inKey(name) {
			value, countValue, err := types.UnmarshalObject(b[count:])
			if err != nil {
				return 0, err
			}
			// value can't be internal values like EOF, Last, Name
			if types.IsInternal(value.Type()) {
				return 0, types.ErrInvalidDocument
			}
			count += countValue
			d.kv[string(name)] = value
		} else {
			_, countValue, err := types.UnmarshalObject(b[count:])
			if err != nil {
				return 0, err
			}
			count += countValue
		}
	}
	return count, nil
}

// used to generate a key that we use in badger for a document
func (d *Document) GetDBKey() []byte {
	if d.key == nil {
		d.key = []byte("document") //FIXME: implement
	}
	return d.key
}

// save this to the database (this is not thread safe)
// func (d *Document) save(tnx *badger.Txn) error {
// 	// save the document to the database
// 	// encode the document
// 	encDoc := types.GenericDocumentUnmarshaler(d)
// 	// save the document to the database
// 	err := tnx.Set(d.GetDBKey(), encDoc)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

type Collection struct {
	name string
	db   *Database

	isSecondary bool
}

// var _ types.Collection = (*Collection)(nil)

func NewCollection(name string, db *Database) *Collection {
	// create a new collection
	return &Collection{
		name: name,
		db:   db,
	}
}

func (c *Collection) Name() string {
	// get the name
	return c.name
}

func (c *Collection) Get(key []byte) (Document, error) {
	panic("not implemented")
}

func (c *Collection) Del(key []byte) error {
	panic("not implemented")
}

func (c *Collection) Set(key []byte, value Document) error {
	panic("not implemented")
}

func (c *Collection) Type() byte {
	return types.CollectionType
}

func (c *Collection) Value() interface{} {
	return c
}

func (c *Collection) String() string {
	return "collection"
}

func (c *Collection) Encode() []byte {
	panic("not implemented")
}

func (c *Collection) Decode([]byte) error {
	panic("not implemented")
}
