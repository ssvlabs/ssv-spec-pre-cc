package messages

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// RoundChangeJustificationsUnmarshalling tests unmarshalling round change justifications
func RoundChangeJustificationsUnmarshalling() tests.SpecTest {

	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], 1, 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[2], 2, 2),
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[3], 3, 2),
	}

	rcMarshalled := testingutils.MarshalJustifications(rcMsgs)

	msg := testingutils.TestingProposalMessageWithParams(
		ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData,
		rcMarshalled, nil)

	msgRoot, err := msg.GetRoot()
	if err != nil {
		panic(err.Error())
	}
	encodedMsg, err := msg.Encode()
	if err != nil {
		panic(err.Error())
	}

	return &tests.MsgSpecTest{
		Name: "round change justification unmarshalling",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		EncodedMessages: [][]byte{
			encodedMsg,
		},
		ExpectedRoots: [][32]byte{
			msgRoot,
		},
	}
}
