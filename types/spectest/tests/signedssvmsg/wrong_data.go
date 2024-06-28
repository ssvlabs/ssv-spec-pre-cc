package signedssvmsg

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// WrongData tests a SignedSSVMessageTest with wrong data (can't decode to SSVMessage)
func WrongData() *SignedSSVMessageTest {
	return &SignedSSVMessageTest{
		Name: "wrong data",
		Messages: []*types.SignedSSVMessage{
			{
				OperatorID: 1,
				Signature:  testingutils.TestingSignedSSVMessageSignature,
				Data:       []byte{1, 2, 3, 4},
			},
		},
		ExpectedError: "could not decode SSVMessage from data in SignedSSVMessage: incorrect size",
	}
}
