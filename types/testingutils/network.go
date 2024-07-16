package testingutils

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
)

type TestingNetwork struct {
	BroadcastedMsgs []*types.SSVMessage
}

func NewTestingNetwork() *TestingNetwork {
	return &TestingNetwork{
		BroadcastedMsgs: make([]*types.SSVMessage, 0),
	}
}

func (net *TestingNetwork) Broadcast(message *types.SSVMessage) error {
	net.BroadcastedMsgs = append(net.BroadcastedMsgs, message)
	return nil
}
