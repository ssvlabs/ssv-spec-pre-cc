package roundchange

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// EmptySigners tests a round change msg with no signers
func EmptySigners() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(5), 2)
	msg.Signers = []types.OperatorID{}

	msgs := []*qbft.SignedMessage{
		msg,
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change empty signer",
		Pre:            pre,
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: invalid signed message: message signers is empty",
	}
}
