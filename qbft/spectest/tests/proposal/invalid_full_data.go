package proposal

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// InvalidFullData tests signed proposal with an invalid full data field (H(full data) != root)
func InvalidFullData() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1))
	msg.FullData = nil

	return &tests.MsgProcessingSpecTest{
		Name:          "invalid full data",
		Pre:           pre,
		PostRoot:      "5b18ca0b470208d8d247543306850618f02bddcbaa7c37eb6d5b36eb3accb5fb",
		InputMessages: []*qbft.SignedMessage{msg},
		ExpectedError: "invalid signed message: H(data) != root",
	}
}
