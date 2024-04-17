package proposal

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// InvalidRoundChangeJustificationPrepared tests a proposal for > 1 round, prepared previously but one of the round change justifications has validRoundChange != nil
// TODO: make sure it does what it used to do before (https://github.com/ssvlabs/ssv-spec-pre-cc/pull/156#discussion_r1113040758)
func InvalidRoundChangeJustificationPrepared() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
		testingutils.TestingPrepareMessage(ks.Shares[2], types.OperatorID(2)),
		testingutils.TestingPrepareMessage(ks.Shares[3], types.OperatorID(3)),
	}
	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithParams(
			ks.Shares[1], types.OperatorID(2), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			qbft.FirstRound, testingutils.MarshalJustifications(prepareMsgs),
		),
		testingutils.TestingRoundChangeMessageWithParams(
			ks.Shares[2], types.OperatorID(2), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			qbft.FirstRound, testingutils.MarshalJustifications(prepareMsgs),
		),
		testingutils.TestingRoundChangeMessageWithParams(
			ks.Shares[3], types.OperatorID(3), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData,
			qbft.FirstRound, testingutils.MarshalJustifications(prepareMsgs),
		),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessageWithParams(ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight,
			testingutils.TestingQBFTRootData,
			testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs),
		),
	}
	return &tests.MsgProcessingSpecTest{
		Name:           "proposal rc msg invalid (prepared)",
		Pre:            pre,
		PostRoot:       "eaa7264b5d6f05cfcdec3158fcae4ff58c3de1e7e9e12bd876177a58686994d4",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: proposal not justified: change round msg not valid: msg signature invalid: failed to verify signature",
	}
}
