package fsblkstorage

import (
	"bytes"
	"context"
	"github.com/hyperledger/fabric/common/ledger"
	coreLedger "github.com/hyperledger/fabric/core/ledger"
	"github.com/hyperledger/fabric/fastfabric/remote"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"sync"
)

func newFsBlockStore(ledgerId string) *BlockStoreImpl {
	return &BlockStoreImpl{ledgerId: ledgerId, client: remote.GetStoragePeerClient(), txCache: sync.Map{}}
}

type BlockStoreImpl struct {
	client        remote.StoragePeerClient
	ledgerId      string
	txCache       sync.Map
	blockHeight   uint64
	currentHash   []byte
	previousHash  []byte
	initialBlocks []*common.Block
}

func (b *BlockStoreImpl) AddBlock(block *common.Block) error {
	if block.Header.Number <= 1 {
		b.initialBlocks = append(b.initialBlocks, block)
	} else {
		if block.Header.Number != b.blockHeight {
			return errors.Errorf(
				"block number should have been %d but was %d",
				b.blockHeight, block.Header.Number,
			)
		}

		if !bytes.Equal(block.Header.PreviousHash, b.currentHash) {
			return errors.Errorf(
				"unexpected Previous block hash. Expected PreviousHash = [%x], PreviousHash referred in the latest block= [%x]",
				b.currentHash, block.Header.PreviousHash,
			)
		}
	}
	b.previousHash = b.currentHash
	b.currentHash = block.Header.Hash()
	b.blockHeight = block.Header.Number + 1
	return nil
}

func (b *BlockStoreImpl) GetBlockchainInfo() (*common.BlockchainInfo, error) {
	return &common.BlockchainInfo{
		Height:            b.blockHeight,
		CurrentBlockHash:  b.currentHash,
		PreviousBlockHash: b.previousHash,
	}, nil
}

type Iterator struct {
	itr    *remote.Iterator
	client remote.StoragePeerClient
}

func (i Iterator) Next() (ledger.QueryResult, error) {
	return i.client.IteratorNext(context.Background(), i.itr)
}

func (i Iterator) Close() {
	_, _ = i.client.IteratorClose(context.Background(), i.itr)
}

func (b *BlockStoreImpl) RetrieveBlocks(startNum uint64) (ledger.ResultsIterator, error) {
	itr, err := b.client.RetrieveBlocks(context.Background(), &remote.RetrieveBlocksRequest{
		LedgerId: b.ledgerId,
		StartNum: startNum})

	return &Iterator{itr: itr, client: b.client}, err
}

func (b *BlockStoreImpl) RetrieveBlockByHash(blockHash []byte) (*common.Block, error) {
	return b.client.RetrieveBlockByHash(context.Background(), &remote.RetrieveBlockByHashRequest{
		LedgerId:  b.ledgerId,
		BlockHash: blockHash})
}

func (b *BlockStoreImpl) RetrieveBlockByNumber(blockNum uint64) (*common.Block, error) {
	if blockNum <= 1 {
		if uint64(len(b.initialBlocks)) < blockNum+1 {
			return nil, nil
		}
		return b.initialBlocks[blockNum], nil
	}
	return b.client.RetrieveBlockByNumber(context.Background(), &remote.RetrieveBlockByNumberRequest{
		LedgerId: b.ledgerId,
		BlockNo:  blockNum})
}

func (b *BlockStoreImpl) checkCache(txID string) bool {
	_, ok := b.txCache.Load(txID)
	return ok
}
func (b *BlockStoreImpl) RetrieveTxByID(txID string) (*common.Envelope, error) {
	if ok := b.checkCache(txID); !ok {
		return nil, coreLedger.NotFoundInIndexErr("")
	}
	tx, err := b.client.RetrieveTxByID(context.Background(), &remote.RetrieveTxByIDRequest{
		LedgerId: b.ledgerId,
		TxID:     txID})
	if err != nil && err.Error() == "rpc error: code = Unknown desc = Entry not found in index" {
		return tx, coreLedger.NotFoundInIndexErr("")
	}
	return tx, err
}

func (b *BlockStoreImpl) RetrieveTxByBlockNumTranNum(blockNum uint64, tranNum uint64) (*common.Envelope, error) {
	return b.client.RetrieveTxByBlockNumTranNum(context.Background(), &remote.RetrieveTxByBlockNumTranNumRequest{
		LedgerId: b.ledgerId,
		BlockNo:  blockNum,
		TxNo:     tranNum})
}

func (b *BlockStoreImpl) RetrieveBlockByTxID(txID string) (*common.Block, error) {
	return b.client.RetrieveBlockByTxID(context.Background(), &remote.RetrieveBlockByTxIDRequest{
		LedgerId: b.ledgerId,
		TxID:     txID})
}

func (b *BlockStoreImpl) RetrieveTxValidationCodeByTxID(txID string) (peer.TxValidationCode, error) {
	code, err := b.client.RetrieveTxValidationCodeByTxID(context.Background(), &remote.RetrieveTxValidationCodeByTxIDRequest{
		LedgerId: b.ledgerId,
		TxID:     txID})
	return peer.TxValidationCode(code.ValidationCode), err
}

func (b *BlockStoreImpl) Shutdown() {
}
