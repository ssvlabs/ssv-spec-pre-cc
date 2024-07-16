package messages

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// InvalidPrepareJustificationsUnmarshalling tests unmarshalling invalid prepare justifications (during message.validate())
func InvalidPrepareJustificationsUnmarshalling() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	msg := testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
		MsgType:              qbft.ProposalMsgType,
		Height:               qbft.FirstHeight,
		Round:                qbft.FirstRound,
		Identifier:           []byte{1, 2, 3, 4},
		Root:                 testingutils.DifferentRoot,
		PrepareJustification: [][]byte{{1}},
	})

	msg.FullData = testingutils.TestingQBFTFullData

	return &tests.MsgSpecTest{
		Name: "invalid prepare justification unmarshalling",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "incorrect size",
	}
}
