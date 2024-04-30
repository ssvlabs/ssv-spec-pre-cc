package messages

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// RoundChangeDataEncoding tests encoding RoundChangeData
func RoundChangeDataEncoding() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingRoundChangeMessageWithParams(
		ks.Shares[1], types.OperatorID(1), qbft.FirstRound, qbft.FirstHeight, testingutils.TestingQBFTRootData, 2,
		testingutils.MarshalJustifications([]*qbft.SignedMessage{
			testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
			testingutils.TestingPrepareMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
			testingutils.TestingPrepareMessageWithRound(ks.Shares[3], types.OperatorID(3), 2),
		}))

	r, _ := msg.GetRoot()
	b, _ := msg.Encode()

	return &tests.MsgSpecTest{
		Name: "round change data encoding",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		EncodedMessages: [][]byte{
			b,
		},
		ExpectedRoots: [][32]byte{
			r,
		},
	}
}
