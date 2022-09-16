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

package fsblkstorage

import (
	"github.com/hyperledger/fabric/common/ledger/blkstorage"
)

// FsBlockstoreProvider provides handle to block storage - this is not thread-safe
type FsBlockstoreProvider struct {
	stores []*BlockStoreImpl
}

// NewProvider constructs a filesystem based block store provider
func NewProvider() blkstorage.BlockStoreProvider {
	return &FsBlockstoreProvider{}
}

// CreateBlockStore simply calls OpenBlockStore
func (p *FsBlockstoreProvider) CreateBlockStore(ledgerid string) (blkstorage.BlockStore, error) {
	return p.OpenBlockStore(ledgerid)
}

// OpenBlockStore opens a block store for given ledgerid.
// If a blockstore is not existing, this method creates one
// This method should be invoked only once for a particular ledgerid
func (p *FsBlockstoreProvider) OpenBlockStore(ledgerid string) (blkstorage.BlockStore, error) {
	newStore := newFsBlockStore(ledgerid)
	p.stores = append(p.stores, newStore)
	return newStore, nil
}

// Exists tells whether the BlockStore with given id exists
func (p *FsBlockstoreProvider) Exists(ledgerid string) (bool, error) {
	for _, store := range p.stores {
		if store.ledgerId == ledgerid {
			return true, nil
		}
	}
	return false, nil
}

// List lists the ids of the existing ledgers
func (p *FsBlockstoreProvider) List() ([]string, error) {
	var ids []string
	for _, store := range p.stores {
		ids = append(ids, store.ledgerId)
	}
	return ids, nil
}

// Close closes the FsBlockstoreProvider
func (p *FsBlockstoreProvider) Close() {
}
