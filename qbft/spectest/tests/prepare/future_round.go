package prepare

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// FutureRound tests prepare for signedProposal.Message.Round > state.Round
func FutureRound() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1))

	msgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessageWithRound(ks.Shares[1], types.OperatorID(1), 3),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "prepare future round",
		Pre:           pre,
		PostRoot:      "b61f5233721865ca43afc68f4ad5045eeb123f6e8f095ce76ecf956dabc74713",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: wrong msg round",
	}
}
