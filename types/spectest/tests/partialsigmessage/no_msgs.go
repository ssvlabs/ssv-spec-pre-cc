package partialsigmessage

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// NoMsgs tests a signed msg with no msgs
func NoMsgs() *MsgSpecTest {
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg.Message.Messages = []*types.PartialSignatureMessage{}

	return &MsgSpecTest{
		Name:          "no messages",
		Messages:      []*types.SignedPartialSignatureMessage{msg},
		ExpectedError: "no PartialSignatureMessages messages",
	}
}
