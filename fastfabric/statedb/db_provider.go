/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package statedb

import (
	"fmt"
	"github.com/hyperledger/fabric/common/ledger/util/leveldbhelper"
	"github.com/hyperledger/fabric/fastfabric/config"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"sync"
)

var logger = logging.MustGetLogger("dbhelper")

var dbNameKeySep = []byte{0x00}
var lastKeyIndicator = byte(0x01)

// Provider enables to use a single leveldb as multiple logical leveldbs
type Provider struct {
	db        *DB
	dbHandles map[string]*DBHandle
	mux       sync.Mutex
	lvlProv   *leveldbhelper.Provider
}

// NewProvider constructs a Provider
func NewProvider(dbPath string) *Provider {
	var p = &Provider{dbHandles: make(map[string]*DBHandle), mux: sync.Mutex{}}
	if config.IsStorage {
		p.lvlProv = leveldbhelper.NewProvider(&leveldbhelper.Conf{DBPath: dbPath})
	} else {
		p.db = createDB()
	}

	return p
}

// GetDBHandle returns a handle to a named db
func (p *Provider) GetDBHandle(dbName string) *DBHandle {
	p.mux.Lock()
	defer p.mux.Unlock()

	if p.dbHandles[dbName] == nil {
		if config.IsStorage {
			fmt.Println("starting leveldb")
			lvlHandle := p.lvlProv.GetDBHandle(dbName)
			p.dbHandles[dbName] = &DBHandle{dbName: dbName, lvlHandle: lvlHandle}
		} else {
			fmt.Println("starting ram db")
			p.dbHandles[dbName] = &DBHandle{dbName: dbName, db: p.db}
		}
	}

	return p.dbHandles[dbName]
}

func (p *Provider) Close() {
	// if config.IsStorage {
	// 	p.lvlProv.Close()
	// } else {
		p.db.Close()
	// }
}

// DBHandle is an handle to a named db
type DBHandle struct {
	dbName    string
	db        *DB
	lvlHandle *leveldbhelper.DBHandle
}

// Get returns the value for the given key
func (h *DBHandle) Get(key []byte) ([]byte, error) {
	// if config.IsStorage {
	// 	return h.lvlHandle.Get(key)
	// } else {
		val, err := h.db.Get(constructLevelKey(h.dbName, key))
		if err == KeyNotFound {
			return nil, nil
		}
		if err != nil {
			logger.Errorf("Error retrieving key [%#v]: %s", key, err)
			return nil, errors.Wrapf(err, "error retrieving key [%#v]", key)
		}
		return val, nil
	// }
}

// Put saves the key/value
func (h *DBHandle) Put(key []byte, value []byte, sync bool) error {
	// if config.IsStorage {
	// 	return h.lvlHandle.Put(key, value, sync)
	// } else {
		return h.db.Put(constructLevelKey(h.dbName, key), value, sync)
	// }
}

// Delete deletes the given key
func (h *DBHandle) Delete(key []byte, sync bool) error {
	// if config.IsStorage {
	// 	return h.lvlHandle.Delete(key, sync)
	// } else {
		return h.db.Delete(constructLevelKey(h.dbName, key), sync)
	// }
}

// WriteBatch writes a batch in an atomic way
func (h *DBHandle) WriteBatch(batch *UpdateBatch, sync bool) error {
	// if config.IsStorage {
	// 	return h.lvlHandle.WriteBatch(batch.lvlBatch, sync)
	// } else {
		for k, v := range batch.KVs {
			key := constructLevelKey(h.dbName, []byte(k))
			if v == nil {
				h.db.Delete(key, true)
			} else {
				h.db.Put(key, v, true)
			}
		}
		return nil
	// }
}

// GetIterator gets an handle to iterator. The iterator should be released after the use.
// The resultset contains all the keys that are present in the db between the startKey (inclusive) and the endKey (exclusive).
// A nil startKey represents the first available key and a nil endKey represent a logical key after the last available key
func (h *DBHandle) GetIterator(startKey []byte, endKey []byte) *Iterator {
	// if config.IsStorage {
	// 	return &Iterator{h.lvlHandle.GetIterator(startKey, endKey)}
	// } else {
		sKey := constructLevelKey(h.dbName, startKey)
		eKey := constructLevelKey(h.dbName, endKey)
		if endKey == nil {
			// replace the last byte 'dbNameKeySep' by 'lastKeyIndicator'
			eKey[len(eKey)-1] = lastKeyIndicator
		}
		logger.Debugf("Getting iterator for range [%#v] - [%#v]", sKey, eKey)
		return &Iterator{Iterator: h.db.GetIterator(sKey, eKey)}
	// }
}

// UpdateBatch encloses the details of multiple `updates`
type UpdateBatch struct {
	KVs      map[string][]byte
	lvlBatch *leveldbhelper.UpdateBatch
}

// NewUpdateBatch constructs an instance of a Batch
func NewUpdateBatch() *UpdateBatch {
	// if config.IsStorage {
	// 	return &UpdateBatch{lvlBatch: leveldbhelper.NewUpdateBatch()}
	// } else {
		return &UpdateBatch{KVs: make(map[string][]byte)}
	// }
}

// Put adds a KV
func (batch *UpdateBatch) Put(key []byte, value []byte) {
	// if config.IsStorage {
	// 	batch.lvlBatch.Put(key, value)
	// } else {
		if value == nil {
			panic("Nil value not allowed")
		}
		batch.KVs[string(key)] = value
	// }
}

// Delete deletes a Key and associated value
func (batch *UpdateBatch) Delete(key []byte) {
	// if config.IsStorage {
	// 	batch.lvlBatch.Delete(key)
	// } else {
		batch.KVs[string(key)] = nil
	// }
}

// Len returns the number of entries in the batch
func (batch *UpdateBatch) Len() int {
	// if config.IsStorage {
	// 	return batch.lvlBatch.Len()
	// } else {
		return len(batch.KVs)
	// }
}

type Iterator struct {
	iterator.Iterator
}

func constructLevelKey(dbName string, key []byte) []byte {
	return append(append([]byte(dbName), dbNameKeySep...), key...)
}
