package roundchange

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// PeerPreparedDifferentHeights tests a round change quorum where peers prepared on different heights
func PeerPreparedDifferentHeights() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 3

	ks := testingutils.Testing4SharesSet()
	prepareMsgs1 := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.Shares[3], types.OperatorID(3)),
	}
	prepareMsgs2 := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 2),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[2], types.OperatorID(2), 2),
		testingutils.TestingPrepareMessageWithRound(ks.Shares[3], types.OperatorID(3), 2),
	}
	msgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 3),
		testingutils.TestingRoundChangeMessageWithParams(
			ks.Shares[2], types.OperatorID(2), 3, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			qbft.FirstRound, testingutils.MarshalJustifications(prepareMsgs1)),
		testingutils.TestingRoundChangeMessageWithParams(ks.Shares[3], types.OperatorID(3), 3, qbft.FirstHeight,
			testingutils.TestingQBFTRootData, 2, testingutils.MarshalJustifications(prepareMsgs2)),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "round change peer prepared different heights",
		Pre:           pre,
		PostRoot:      "802b2c5b3e85312a7068a28a22d91f0af050c2bf17993730896929cd40f50195",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), 3, qbft.FirstHeight,
				testingutils.TestingQBFTRootData,
				testingutils.MarshalJustifications(msgs), testingutils.MarshalJustifications(prepareMsgs2)),
		},
	}
}
