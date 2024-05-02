package partialsigmessage

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// InconsistentSignedMessage tests SignedPartialSignatureMessage where the signer is not the same as the signer in messages
func InconsistentSignedMessage() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb)
	msgWithDifferentSigner := testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, qbft.FirstHeight)

	msg.Message.Messages = append(msg.Message.Messages, msgWithDifferentSigner.Message.Messages...)

	return &MsgSpecTest{
		Name: "inconsistent signed message",
		Messages: []*types.SignedPartialSignatureMessage{
			msg,
		},
		ExpectedError: "inconsistent signers",
	}
}
