package signedssvmsg

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// ZeroSigner tests an invalid SignedSSVMessageTest with zero signer
func ZeroSigner() *SignedSSVMessageTest {

	return &SignedSSVMessageTest{
		Name: "zero signer",
		Messages: []*types.SignedSSVMessage{
			{
				OperatorID: 0,
				Signature:  testingutils.TestingSignedSSVMessageSignature,
				Data:       []byte{1, 2, 3, 4},
			},
		},
		ExpectedError: "OperatorID in SignedSSVMessage is 0",
	}
}
