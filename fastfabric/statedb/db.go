package statedb

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/common/ledger/util/leveldbhelper"
	// "github.com/hyperledger/fabric/fastfabric/config"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var KeyNotFound = errors.New("key not found")

type DB struct {
	db    *ValueHashtable
	lvldb *leveldbhelper.DB
}

func createDB() *DB {
	db := &DB{db: NewHT()}
	return db
}

// func (ledger *DB) Open() {
// 	if config.IsStorage {
// 		ledger.lvldb.Open()
// 	}
// }

func (ledger *DB) Close() error {
	// if config.IsStorage {
	// 	ledger.lvldb.Close()
	// 	return nil
	// }
	ledger.db.Cleanup()
	return nil
}

func (ledger *DB) Get(bytes []byte) ([]byte, error) {
	// if config.IsStorage {
	// 	return ledger.lvldb.Get(bytes)
	// }

	if val, err := ledger.db.Get(bytes); err != nil {
		return nil, KeyNotFound
	} else {
		return val, nil
	}
}
func (ledger *DB) Put(key []byte, value []byte, sync bool) error {
	// if config.IsStorage {
	// 	return ledger.lvldb.Put(key, value, sync)
	// }

	return ledger.db.Put(key, value)
}
func (ledger *DB) Delete(key []byte, sync bool) error {
	// if config.IsStorage {
	// 	return ledger.lvldb.Delete(key, sync)
	// }

	return ledger.db.Remove(key)
}
func (ledger *DB) GetIterator(sk []byte, ek []byte) iterator.Iterator {
	// if config.IsStorage {
	// 	return ledger.lvldb.GetIterator(sk, ek)
	// }

	keys := ledger.db.getKeys(sk, ek)
	itr := &Iterimpl{keys: keys, idx: -1, values: make(map[string][]byte)}
	for _, key := range keys {
		v, err := ledger.db.Get(key)
		if err != nil {
			itr.err = fmt.Errorf("Failed to get value for key %s: %s", key, err)
			return itr
		}
		itr.values[string(retrieveAppKey(key))] = v
	}
	return itr
}

type Iterimpl struct {
	keys   [][]byte
	idx    int
	err    error
	values map[string][]byte
}

func (i *Iterimpl) First() bool {
	i.idx = 0
	return i.keys != nil
}

func (i *Iterimpl) Last() bool {
	i.idx = len(i.keys) - 1
	return i.keys != nil
}

func (i *Iterimpl) Seek(key []byte) bool {
	for tempidx := 0; tempidx < len(i.keys); tempidx++ {
		if bytes.Compare(key, i.keys[tempidx]) > -1 {
			i.idx = tempidx
			return true
		}
	}
	return false
}

func (i *Iterimpl) Next() bool {
	i.idx++
	return i.idx < len(i.keys)
}

func (i *Iterimpl) Prev() bool {
	i.idx--
	return i.idx >= 0
}

func (Iterimpl) Release() {
	return
}

func (Iterimpl) SetReleaser(releaser util.Releaser) {
	return
}

func (i *Iterimpl) Valid() bool {
	return i.idx >= 0 && i.idx < len(i.keys)
}

func (i *Iterimpl) Error() error {
	return i.err
}

func (i *Iterimpl) Key() []byte {
	if !i.Valid() {
		return nil
	}
	return retrieveAppKey(i.keys[i.idx])
}

func retrieveAppKey(levelKey []byte) []byte {
	return bytes.SplitN(levelKey, dbNameKeySep, 2)[1]
}

func (i *Iterimpl) Value() []byte {
	if !i.Valid() {
		return nil
	}

	return i.values[string(i.Key())]
}
