package roundchange

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// PostCutoff tests processing a round change msg when round >= cutoff
func PostCutoff() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	pre := testingutils.BaseInstance()
	pre.State.Round = 15

	msgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], types.OperatorID(1), 16),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "round cutoff round change message",
		Pre:           pre,
		PostRoot:      "9c9b1da0c431638ec4bbaabd98dfdca67ec54f6b56b2162ff7faf749c4efdcab",
		InputMessages: msgs,
		ExpectedError: "instance stopped processing messages",
	}
}
