package messages

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// RoundChangePrePreparedJustifications tests valid justified change round (prev prepared)
func RoundChangePrePreparedJustifications() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.Shares[3], types.OperatorID(3)),
	}

	msg := testingutils.TestingRoundChangeMessageWithParams(
		ks.Shares[1], types.OperatorID(1), 10, qbft.FirstHeight, testingutils.TestingQBFTRootData,
		qbft.FirstRound, testingutils.MarshalJustifications(prepareMsgs))

	return &tests.MsgSpecTest{
		Name: "rc prev prepared justifications",
		Messages: []*qbft.SignedMessage{
			msg,
		},
	}
}
