package statedb

import (
	"bytes"
	"errors"
	"sync"
)

// ValueHashtable the set of Items
type ValueHashtable struct {
	items map[string][]byte
	lock  sync.RWMutex
}

func NewHT() *ValueHashtable {
	return &ValueHashtable{lock: sync.RWMutex{}}
}

func (ht *ValueHashtable) Put(key []byte, value []byte) error {
	ht.lock.Lock()
	defer ht.lock.Unlock()
	if ht.items == nil {
		ht.items = make(map[string][]byte)
	}
	ht.items[string(key)] = value
	return nil
}

// Remove item with key k from hashtable
func (ht *ValueHashtable) Remove(key []byte) error {
	ht.lock.Lock()
	defer ht.lock.Unlock()
	delete(ht.items, string(key))
	return nil
}

// Get item with key k from the hashtable
func (ht *ValueHashtable) Get(key []byte) ([]byte, error) {
	ht.lock.RLock()
	defer ht.lock.RUnlock()
	if val, ok := ht.items[string(key)]; ok {
		return val, nil
	} else {
		return nil, errors.New("key not found")
	}
}

// Size returns the number of the hashtable elements
func (ht *ValueHashtable) Size() int {
	ht.lock.RLock()
	defer ht.lock.RUnlock()
	return len(ht.items)
}

func (ht *ValueHashtable) Cleanup() {
	ht.lock.RLock()
	defer ht.lock.RUnlock()
	ht.items = nil
}

func (ht *ValueHashtable) getKeys(sk []byte, ek []byte) [][]byte {
	ht.lock.RLock()
	defer ht.lock.RUnlock()
	keys := make([][]byte, 0)

	for k := range ht.items {
		x := []byte(k)
		if bytes.Compare(sk, x) < 1 && bytes.Compare(x, ek) < 1 {
			keys = append(keys, x)
		}
	}
	return keys
}
