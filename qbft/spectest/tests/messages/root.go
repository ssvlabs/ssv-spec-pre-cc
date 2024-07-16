package messages

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// GetRoot tests GetRoot on SignedMessage
func GetRoot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingProposalMessageWithParams(
		ks.Shares[1], types.OperatorID(1), qbft.FirstRound, qbft.FirstHeight, testingutils.TestingQBFTRootData,
		testingutils.MarshalJustifications([]*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		}),
		testingutils.MarshalJustifications([]*qbft.SignedMessage{
			testingutils.TestingRoundChangeMessage(ks.Shares[1], types.OperatorID(1)),
		}))

	r, _ := msg.GetRoot()

	return &tests.MsgSpecTest{
		Name: "get root",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedRoots: [][32]byte{
			r,
		},
	}
}
