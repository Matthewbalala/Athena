package cached

import (
	"fmt"
	"github.com/hyperledger/fabric/protos/common"
)

func (b *Block) GetChannelId() (string, error) {
	env, err := b.UnmarshalSpecificEnvelope(0)
	if err != nil {
		return "", err
	}
	return env.GetChannelId()
}

func (env *Envelope) GetChannelId() (string, error) {
	pl, err := env.UnmarshalPayload()
	if err != nil {
		return "", err
	}
	hdr, err := pl.Header.UnmarshalChannelHeader()
	if err != nil {
		return "", err
	}
	return hdr.ChannelId, nil
}

// IsConfigBlock validates whenever given block contains configuration
// update transaction
func (block *Block) IsConfigBlock() bool {
	envelope, err := block.UnmarshalSpecificEnvelope(0)
	if err != nil {
		return false
	}

	payload, err := envelope.UnmarshalPayload()
	if err != nil {
		return false
	}

	if payload.Header == nil {
		return false
	}

	hdr, err := payload.Header.UnmarshalChannelHeader()
	if err != nil {
		return false
	}

	return common.HeaderType(hdr.Type) == common.HeaderType_CONFIG ||
		common.HeaderType(hdr.Type) == common.HeaderType_ORDERER_TRANSACTION
}

func (env *Envelope) GetChaincodeAction() (*ChaincodeAction, error) {
	pl, err := env.UnmarshalPayload()
	if err != nil {
		return nil, err
	}
	tx, err := pl.UnmarshalTransaction()
	if err != nil {
		return nil, err
	}
	if len(tx.Actions) == 0 {
		return nil, fmt.Errorf("no actions in transaction")
	}
	cap, err := tx.Actions[0].UnmarshalChaincodeActionPayload()
	if err != nil {
		return nil, err
	}
	prp, err := cap.Action.UnmarshalProposalResponsePayload()
	if err != nil {
		return nil, err
	}
	return prp.UnmarshalChaincodeAction()
}
