package processmsg

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// MsgError tests a process msg returning an error
func MsgError() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	return &tests.ControllerSpecTest{
		Name: "process msg error",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.TestingProposalMessageWithRound(ks.Shares[1], 1, 100),
				},
				ControllerPostRoot: "47713c38fe74ce55959980781287886c603c2117a14dc8abce24dcb9be0093af",
			},
		},
		ExpectedError: "could not process msg: invalid signed message: proposal not justified: change round has no quorum",
	}
}
