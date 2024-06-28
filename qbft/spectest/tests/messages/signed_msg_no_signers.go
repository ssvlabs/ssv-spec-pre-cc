package messages

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// SignedMsgNoSigners tests SignedMessage len(signers) == 0
func SignedMsgNoSigners() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingCommitMessage(ks.Shares[1], types.OperatorID(1))
	msg.Signers = nil

	return &tests.MsgSpecTest{
		Name: "no signers",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "message signers is empty",
	}
}
