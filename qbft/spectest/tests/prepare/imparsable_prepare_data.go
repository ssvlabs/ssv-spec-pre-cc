package prepare

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ImparsableProposalData tests a prepare msg received with imparsable data
func ImparsableProposalData() *tests.MsgProcessingSpecTest {
	pre := testingutils.BaseInstance()
	pre.State.ProposalAcceptedForCurrentRound = testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
	})
	msgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(testingutils.Testing4SharesSet().Shares[3], types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       []byte{1, 2, 3, 4},
		}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "imparsable prepare data",
		Pre:           pre,
		PostRoot:      "be41977d818071451988105377df7c5ccf89ecc05ddf033b7b3b83d89f52d530",
		InputMessages: msgs,
		ExpectedError: "invalid signed message: could not get prepare data: could not decode prepare data from message: invalid character '\\x01' looking for beginning of value",
	}
}
