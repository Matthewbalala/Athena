package cached

import (
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/peer"
	"sync"
)

type Block struct {
	*common.Block
	cachedEnvs     []*Envelope
	cachedMetadata []*Metadata
	metaMtx        sync.RWMutex
	envMtx         sync.RWMutex
}

type Metadata struct {
	*common.Metadata
	cachedSigHeaders []*common.SignatureHeader
	sigMtx           sync.Mutex
}

type Envelope struct {
	*common.Envelope
	cachedPayload *Payload
	plMtx         sync.Mutex
}

type Payload struct {
	*common.Payload
	cachedEnTx *Transaction
	Header     *Header
	txMtx      sync.Mutex
}

type ChannelHeader struct {
	*common.ChannelHeader
	cachedExtension *peer.ChaincodeHeaderExtension
	extMtx          sync.Mutex
}

type Transaction struct {
	*peer.Transaction
	Actions []*TransactionAction
}

type TransactionAction struct {
	*peer.TransactionAction
	cachedSigHeader     *common.SignatureHeader
	cachedActionPayload *ChaincodeActionPayload
	sigMtx              sync.Mutex
	plMtx               sync.Mutex
}

type ProposalResponsePayload struct {
	*peer.ProposalResponsePayload
	cachedAction *ChaincodeAction
	actMtx       sync.Mutex
}

type ChaincodeInvocationSpec struct {
	*peer.ChaincodeInvocationSpec
}

type ChaincodeProposalPayload struct {
	*peer.ChaincodeProposalPayload
	cachedInput *ChaincodeInvocationSpec
	inMtx       sync.Mutex
}

type ChaincodeActionPayload struct {
	*peer.ChaincodeActionPayload
	Action            *ChaincodeEndorsedAction
	cachedPropPayload *ChaincodeProposalPayload
	propMtx           sync.Mutex
}

type ChaincodeEndorsedAction struct {
	*peer.ChaincodeEndorsedAction
	cachedRespPayload *ProposalResponsePayload
	propMtx           sync.Mutex
}

type ChaincodeAction struct {
	*peer.ChaincodeAction
	cachedRwSet  *TxRwSet
	cachedEvents *peer.ChaincodeEvent
	setMtx       sync.Mutex
	evMtx        sync.Mutex
}

type Header struct {
	*common.Header
	cachedChanHeader *ChannelHeader
	cachedSigHeader  *common.SignatureHeader
	chdrMtx          sync.Mutex
	sigMtx           sync.Mutex
}

type GossipPayload struct {
	Data        *Block `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	PrivateData [][]byte
}
