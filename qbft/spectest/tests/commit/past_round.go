package commit

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// PastRound tests a commit msg with past round, should process but not decide
func PastRound() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessageWithRound(ks.Shares[1], 1, 5)
	pre.State.Round = 5

	msgs := []*qbft.SignedMessage{
		testingutils.TestingCommitMessageWithRound(ks.Shares[1], 1, 2),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "commit past round",
		Pre:           pre,
		PostRoot:      "929e208f33ae4f872976a962256ddddebfb625905ff02903315586688ec1562d",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: past round",
	}
}
