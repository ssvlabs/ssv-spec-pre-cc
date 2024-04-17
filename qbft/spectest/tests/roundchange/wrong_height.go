package roundchange

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// WrongHeight tests a round change msg with wrong height
func WrongHeight() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRoundAndHeight(ks.Shares[1], types.OperatorID(1), 2, 2),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change invalid height",
		Pre:            pre,
		PostRoot:       "96e6d7bdbb98a2d9937f3d97d6aa096bd3a58f923b61012048ac70ad52765919",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: wrong msg height",
	}
}
