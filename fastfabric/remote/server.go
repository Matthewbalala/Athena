package remote

import (
	"context"
	"errors"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/common/ledger"
	"github.com/hyperledger/fabric/common/ledger/blkstorage"
	"github.com/hyperledger/fabric/fastfabric/config"
	"github.com/hyperledger/fabric/protos/common"
	"google.golang.org/grpc"
	"net"
)

var remoteLogger = flogging.MustGetLogger("remote")
var StorageServer = &server{stores: make(map[string]blkstorage.BlockStore), err: make(chan error, 1)}

func StartServer(address string) {
	config.RegisterBlockStore = func(id string, store interface{}) {
		StorageServer.RegisterBlockStore(id, store.(blkstorage.BlockStore))
	}
	lis, err := net.Listen("tcp", address)
	if err != nil {
		remoteLogger.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterStoragePeerServer(s, StorageServer)
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

type server struct {
	iterators []ledger.ResultsIterator
	err       chan error
	stores    map[string]blkstorage.BlockStore
}

func (s *server) RegisterBlockStore(ledgerId string, store blkstorage.BlockStore) {
	s.stores[ledgerId] = store
}

func (s *server) IteratorNext(ctx context.Context, itr *Iterator) (*common.Block, error) {
	res, err := s.iterators[itr.IteratorId].Next()
	if res == nil {
		return nil, err
	}
	return res.(*common.Block), err
}

func (s *server) IteratorClose(ctx context.Context, itr *Iterator) (*Result, error) {
	s.iterators[itr.IteratorId].Close()
	s.iterators[itr.IteratorId] = nil
	return &Result{}, nil
}

func (s *server) RetrieveBlocks(ctx context.Context, req *RetrieveBlocksRequest) (*Iterator, error) {
	if store, ok := s.stores[req.LedgerId]; ok {
		itr, err := store.RetrieveBlocks(req.StartNum)
		s.iterators = append(s.iterators, itr)
		return &Iterator{IteratorId: int32(len(s.iterators) - 1)}, err
	}
	return nil, errors.New("store not initialized yet")
}

func (s *server) RetrieveTxValidationCodeByTxID(ctx context.Context, req *RetrieveTxValidationCodeByTxIDRequest) (*ValidationCode, error) {
	if store, ok := s.stores[req.LedgerId]; ok {
		code, err := store.RetrieveTxValidationCodeByTxID(req.TxID)
		return &ValidationCode{ValidationCode: int32(code)}, err
	}
	return nil, errors.New("store not initialized yet")
}

func (s *server) RetrieveBlockByTxID(ctx context.Context, req *RetrieveBlockByTxIDRequest) (*common.Block, error) {
	if store, ok := s.stores[req.LedgerId]; ok {
		return store.RetrieveBlockByTxID(req.TxID)
	}
	return nil, errors.New("store not initialized yet")
}

func (s *server) RetrieveTxByBlockNumTranNum(ctx context.Context, req *RetrieveTxByBlockNumTranNumRequest) (*common.Envelope, error) {
	if store, ok := s.stores[req.LedgerId]; ok {
		return store.RetrieveTxByBlockNumTranNum(req.BlockNo, req.TxNo)
	}
	return nil, errors.New("store not initialized yet")
}

func (s *server) RetrieveTxByID(ctx context.Context, req *RetrieveTxByIDRequest) (*common.Envelope, error) {
	if store, ok := s.stores[req.LedgerId]; ok {
		return store.RetrieveTxByID(req.TxID)
	}
	return nil, errors.New("store not initialized yet")
}

func (s *server) RetrieveBlockByNumber(ctx context.Context, req *RetrieveBlockByNumberRequest) (*common.Block, error) {
	if store, ok := s.stores[req.LedgerId]; ok {
		return store.RetrieveBlockByNumber(req.BlockNo)
	}
	return nil, errors.New("store not initialized yet")
}

func (s *server) GetBlockchainInfo(ctx context.Context, req *GetBlockchainInfoRequest) (*common.BlockchainInfo, error) {
	if store, ok := s.stores[req.LedgerId]; ok {
		return store.GetBlockchainInfo()
	}
	return &common.BlockchainInfo{
		Height:            0,
		CurrentBlockHash:  nil,
		PreviousBlockHash: nil}, nil
}

func (s *server) RetrieveBlockByHash(ctx context.Context, req *RetrieveBlockByHashRequest) (*common.Block, error) {
	if store, ok := s.stores[req.LedgerId]; ok {
		return store.RetrieveBlockByHash(req.BlockHash)
	}
	return nil, errors.New("store not initialized yet")
}
