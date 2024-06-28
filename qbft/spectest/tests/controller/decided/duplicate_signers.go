package decided

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// DuplicateSigners tests a decided msg with duplicate signers
func DuplicateSigners() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingCommitMultiSignerMessageWithHeight([]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]}, []types.OperatorID{1, 2, 3}, 10)
	msg.Signers = []types.OperatorID{1, 2, 2}

	return &tests.ControllerSpecTest{
		Name: "decide duplicate signer",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					msg,
				},
				ControllerPostRoot: "47713c38fe74ce55959980781287886c603c2117a14dc8abce24dcb9be0093af",
			},
		},
		ExpectedError: "invalid decided msg: invalid decided msg: non unique signer",
	}
}
