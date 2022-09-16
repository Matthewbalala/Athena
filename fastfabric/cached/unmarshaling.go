package cached

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
)

func WrapBlock(raw *common.Block) *Block {
	if raw == nil {
		return nil
	}

	lenMeta := 0
	if raw.Metadata != nil {
		lenMeta = len(raw.Metadata.Metadata)
	}

	lenEnvs := 0
	if raw.Data != nil {
		lenEnvs = len(raw.Data.Data)
	}

	return &Block{
		Block:          raw,
		cachedMetadata: make([]*Metadata, lenMeta),
		cachedEnvs:     make([]*Envelope, lenEnvs)}
}

func (b *Block) UnmarshalAll() error {
	metas, err := b.UnmarshalAllMetadata()
	if err != nil {
		return err
	}
	for _, meta := range metas {
		if _, err := meta.UnmarshalAllSignatureHeaders(); err != nil {
			return err
		}
	}
	envs, err := b.UnmarshalAllEnvelopes()
	if err != nil {
		return err
	}
	for _, env := range envs {
		pl, err := env.UnmarshalPayload()
		if err != nil {
			return err
		}
		chdr, err := pl.Header.UnmarshalChannelHeader()
		if err != nil {
			return err
		}
		_, err = chdr.unmarshalExtension()
		if err != nil {
			return err
		}

		_, err = pl.Header.UnmarshalSignatureHeader()
		if err != nil {
			return err
		}

		etx, err := pl.UnmarshalTransaction()
		if err != nil {
			return err
		}

		for _, act := range etx.Actions {
			_, err := act.UnmarshalSignatureHeader()
			if err != nil {
				return err
			}
			pl, err := act.UnmarshalChaincodeActionPayload()
			if err != nil {
				return err
			}

			respPl, err := pl.Action.UnmarshalProposalResponsePayload()
			if err != nil {
				return err
			}

			act, err := respPl.UnmarshalChaincodeAction()
			if err != nil {
				return err
			}

			_, err = act.UnmarshalRwSet()
			if err != nil {
				return err
			}

			_, err = act.UnmarshalEvents()
			if err != nil {
				return err
			}

			propPl, err := pl.UnmarshalProposalPayload()
			if err != nil {
				return err
			}

			_, err = propPl.UnmarshalInput()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Block) UnmarshalAllMetadata() ([]*Metadata, error) {
	if b.Metadata == nil {
		return nil, fmt.Errorf("block metadata must not be nil")
	}

	for i := range b.Metadata.Metadata {
		if _, err := b.UnmarshalSpecificMetadata(common.BlockMetadataIndex(i)); err != nil {
			return nil, err
		}
	}
	return b.cachedMetadata, nil
}

func (b *Block) UnmarshalSpecificMetadata(index common.BlockMetadataIndex) (*Metadata, error) {
	b.metaMtx.Lock()
	defer b.metaMtx.Unlock()

	if len(b.cachedMetadata) <= int(index) || index < 0 {
		return nil, fmt.Errorf("index out of range")
	}

	if b.cachedMetadata[index] != nil {
		return b.cachedMetadata[index], nil
	}

	metaRaw := &common.Metadata{}
	if err := proto.Unmarshal(b.Metadata.Metadata[index], metaRaw); err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling metadata from block at index [%s]", index)
	}

	meta := &Metadata{
		Metadata:         metaRaw,
		cachedSigHeaders: make([]*common.SignatureHeader, len(metaRaw.Signatures))}
	b.cachedMetadata[index] = meta
	return meta, nil
}

func (meta *Metadata) UnmarshalSpecificSignatureHeader(index int) (*common.SignatureHeader, error) {
	meta.sigMtx.Lock()
	defer meta.sigMtx.Unlock()
	if len(meta.cachedSigHeaders) <= int(index) || index < 0 {
		return nil, fmt.Errorf("index out of range")
	}
	if meta.cachedSigHeaders[index] != nil {
		return meta.cachedSigHeaders[index], nil
	}

	var err error
	meta.cachedSigHeaders[index], err = unmarshalSignatureHeader(meta.Signatures[index].SignatureHeader)
	return meta.cachedSigHeaders[index], err
}

func (meta *Metadata) UnmarshalAllSignatureHeaders() ([]*common.SignatureHeader, error) {
	for i := range meta.Signatures {
		_, err := meta.UnmarshalSpecificSignatureHeader(i)
		if err != nil {
			return nil, err
		}
	}
	return meta.cachedSigHeaders, nil
}

func (b *Block) UnmarshalAllEnvelopes() ([]*Envelope, error) {
	if b.Data == nil || b.Data.Data == nil {
		return nil, fmt.Errorf("block data must not be nil")
	}

	for i := range b.Data.Data {
		if _, err := b.UnmarshalSpecificEnvelope(i); err != nil {
			return nil, err
		}
	}

	return b.cachedEnvs, nil
}

func (b *Block) UnmarshalSpecificEnvelope(index int) (*Envelope, error) {
	b.envMtx.Lock()
	defer b.envMtx.Unlock()
	if b.Data == nil || b.Data.Data == nil {
		return nil, fmt.Errorf("block data must not be nil")
	}
	if len(b.cachedEnvs) <= int(index) || index < 0 {
		return nil, fmt.Errorf("index out of range")
	}
	if b.cachedEnvs[index] != nil {
		return b.cachedEnvs[index], nil
	}

	envRaw := &common.Envelope{}
	if err := proto.Unmarshal(b.Data.Data[index], envRaw); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling Envelope")
	}

	env := &Envelope{Envelope: envRaw}
	b.cachedEnvs[index] = env
	return env, nil
}

func (env *Envelope) UnmarshalPayload() (*Payload, error) {
	env.plMtx.Lock()
	defer env.plMtx.Unlock()
	if env.cachedPayload != nil {
		return env.cachedPayload, nil
	}

	payloadRaw := &common.Payload{}
	if err := proto.Unmarshal(env.Payload, payloadRaw); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling Payload")
	}

	payload := &Payload{Payload: payloadRaw, Header: &Header{Header: payloadRaw.Header}}
	env.cachedPayload = payload
	return payload, nil
}

func (hdr *Header) UnmarshalChannelHeader() (*ChannelHeader, error) {
	hdr.chdrMtx.Lock()
	defer hdr.chdrMtx.Unlock()
	if hdr.cachedChanHeader != nil {
		return hdr.cachedChanHeader, nil
	}

	if hdr.Header == nil {
		return nil, fmt.Errorf("payload header is nil")
	}

	headerRaw := &common.ChannelHeader{}
	if err := proto.Unmarshal(hdr.Header.ChannelHeader, headerRaw); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling payload ChannelHeader")
	}

	header := &ChannelHeader{ChannelHeader: headerRaw}
	hdr.cachedChanHeader = header
	return header, nil
}

func (hdr *Header) UnmarshalSignatureHeader() (*common.SignatureHeader, error) {
	hdr.sigMtx.Lock()
	defer hdr.sigMtx.Unlock()
	if hdr.cachedSigHeader != nil {
		return hdr.cachedSigHeader, nil
	}

	if hdr.Header == nil {
		return nil, fmt.Errorf("payload header is nil")
	}

	headerRaw, err := unmarshalSignatureHeader(hdr.Header.SignatureHeader)

	hdr.cachedSigHeader = headerRaw
	return headerRaw, err
}

func (hdr *Header) UnmarshalChaincodeHeaderExtension() (*peer.ChaincodeHeaderExtension, error) {
	chdr, err := hdr.UnmarshalChannelHeader()
	if err != nil {
		return nil, err
	}
	return chdr.unmarshalExtension()
}

func unmarshalSignatureHeader(bytes []byte) (*common.SignatureHeader, error) {
	headerRaw := &common.SignatureHeader{}
	if err := proto.Unmarshal(bytes, headerRaw); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling SignatureHeader")
	}
	return headerRaw, nil
}

func (ch *ChannelHeader) unmarshalExtension() (*peer.ChaincodeHeaderExtension, error) {
	ch.extMtx.Lock()
	defer ch.extMtx.Unlock()
	if ch.cachedExtension != nil {
		return ch.cachedExtension, nil
	}

	ext := &peer.ChaincodeHeaderExtension{}
	if err := proto.Unmarshal(ch.Extension, ext); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling channel header ChaincodeHeaderExtension")
	}
	ch.cachedExtension = ext
	return ext, nil
}

func (pl *Payload) UnmarshalTransaction() (*Transaction, error) {
	pl.txMtx.Lock()
	defer pl.txMtx.Unlock()
	if pl.cachedEnTx != nil {
		return pl.cachedEnTx, nil
	}
	txRaw := &peer.Transaction{}
	if err := proto.Unmarshal(pl.Data, txRaw); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling Transaction")
	}

	tx := &Transaction{Transaction: txRaw, Actions: make([]*TransactionAction, len(txRaw.Actions))}
	for i, a := range txRaw.Actions {
		tx.Actions[i] = &TransactionAction{TransactionAction: a}
	}
	pl.cachedEnTx = tx
	return tx, nil
}

func (pl *Payload) UnmarshalChaincodeAction() (*ChaincodeAction, error) {
	tx, err := pl.UnmarshalTransaction()
	if err != nil {
		return nil, err
	}
	capl, err := tx.UnmarshalChaincodeActionPayload()
	if err != nil {
		return nil, err
	}
	prpl, err := capl.Action.UnmarshalProposalResponsePayload()
	if err != nil {
		return nil, err
	}
	return prpl.UnmarshalChaincodeAction()

}

func (tx *Transaction) UnmarshalChaincodeActionPayload() (*ChaincodeActionPayload, error) {
	if len(tx.Actions) < 1 {
		return nil, errors.New("transaction should have at least one action")
	}
	capl, err := tx.Actions[0].UnmarshalChaincodeActionPayload()
	if err != nil {
		return nil, err
	}
	return capl, nil
}

func (act *TransactionAction) UnmarshalSignatureHeader() (*common.SignatureHeader, error) {
	act.sigMtx.Lock()
	defer act.sigMtx.Unlock()
	if act.cachedSigHeader != nil {
		return act.cachedSigHeader, nil
	}

	var err error
	act.cachedSigHeader, err = unmarshalSignatureHeader(act.Header)
	return act.cachedSigHeader, err
}

func (act *TransactionAction) UnmarshalChaincodeActionPayload() (*ChaincodeActionPayload, error) {
	act.plMtx.Lock()
	defer act.plMtx.Unlock()
	if act.cachedActionPayload != nil {
		return act.cachedActionPayload, nil
	}

	capRaw := &peer.ChaincodeActionPayload{}
	if err := proto.Unmarshal(act.Payload, capRaw); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling ChaincodeActionPayload")
	}

	act.cachedActionPayload = &ChaincodeActionPayload{
		ChaincodeActionPayload: capRaw,
		Action: &ChaincodeEndorsedAction{
			ChaincodeEndorsedAction: capRaw.Action}}
	return act.cachedActionPayload, nil
}

func (pl *ChaincodeActionPayload) UnmarshalProposalPayload() (*ChaincodeProposalPayload, error) {
	pl.propMtx.Lock()
	defer pl.propMtx.Unlock()
	if pl.cachedPropPayload != nil {
		return pl.cachedPropPayload, nil
	}

	cpp := &peer.ChaincodeProposalPayload{}
	if err := proto.Unmarshal(pl.ChaincodeProposalPayload, cpp); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling ChaincodeProposalPayload")
	}

	pl.cachedPropPayload = &ChaincodeProposalPayload{ChaincodeProposalPayload: cpp}
	return pl.cachedPropPayload, nil
}

func (act *ChaincodeEndorsedAction) UnmarshalProposalResponsePayload() (*ProposalResponsePayload, error) {
	act.propMtx.Lock()
	defer act.propMtx.Unlock()
	if act.cachedRespPayload != nil {
		return act.cachedRespPayload, nil
	}

	prp := &peer.ProposalResponsePayload{}
	if err := proto.Unmarshal(act.ProposalResponsePayload, prp); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling ProposalResponsePayload")
	}

	act.cachedRespPayload = &ProposalResponsePayload{ProposalResponsePayload: prp}
	return act.cachedRespPayload, nil
}

func (respPl *ProposalResponsePayload) UnmarshalChaincodeAction() (*ChaincodeAction, error) {
	respPl.actMtx.Lock()
	defer respPl.actMtx.Unlock()
	if respPl.cachedAction != nil {
		return respPl.cachedAction, nil
	}

	chaincodeAction := &peer.ChaincodeAction{}
	if err := proto.Unmarshal(respPl.Extension, chaincodeAction); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling ChaincodeAction")
	}

	respPl.cachedAction = &ChaincodeAction{ChaincodeAction: chaincodeAction}
	return respPl.cachedAction, nil
}

func (act *ChaincodeAction) UnmarshalRwSet() (*TxRwSet, error) {
	act.setMtx.Lock()
	defer act.setMtx.Unlock()

	if act.cachedRwSet != nil {
		return act.cachedRwSet, nil
	}

	rwset := &TxRwSet{}
	if err := rwset.FromProtoBytes(act.Results); err != nil {
		return nil, err
	}

	act.cachedRwSet = rwset
	return rwset, nil
}

func (act *ChaincodeAction) UnmarshalEvents() (*peer.ChaincodeEvent, error) {
	act.evMtx.Lock()
	defer act.evMtx.Unlock()
	if act.cachedEvents != nil {
		return act.cachedEvents, nil
	}
	event := &peer.ChaincodeEvent{}
	if err := proto.Unmarshal(act.Events, event); err != nil {
		return nil, err
	}

	act.cachedEvents = event
	return event, nil
}

func (cpp *ChaincodeProposalPayload) UnmarshalInput() (*ChaincodeInvocationSpec, error) {
	cpp.inMtx.Lock()
	cpp.inMtx.Unlock()
	if cpp.cachedInput != nil {
		return cpp.cachedInput, nil
	}

	cis := &peer.ChaincodeInvocationSpec{}
	if err := proto.Unmarshal(cpp.Input, cis); err != nil {
		return nil, err
	}
	cpp.cachedInput = &ChaincodeInvocationSpec{ChaincodeInvocationSpec: cis}
	return cpp.cachedInput, nil
}
