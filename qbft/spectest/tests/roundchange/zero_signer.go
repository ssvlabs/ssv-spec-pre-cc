package roundchange

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// ZeroSigner tests a signed round change msg with signer 0
func ZeroSigner() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(0), 2),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change zero signer",
		Pre:            pre,
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: invalid signed message: signer ID 0 not allowed",
	}
}
