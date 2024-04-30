package roundchange

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// NonUniqueSigner tests a round change msg with multiple signers and non unique signer
func NonUniqueSigner() tests.SpecTest {
	pre := testingutils.BaseInstance()
	pre.State.Round = 2
	ks := testingutils.Testing4SharesSet()

	msg := testingutils.TestingMultiSignerRoundChangeMessageWithRound(
		[]*bls.SecretKey{ks.Shares[1], ks.Shares[2]},
		[]types.OperatorID{types.OperatorID(1), types.OperatorID(2)},
		2,
	)
	msg.Signers = []types.OperatorID{types.OperatorID(1), types.OperatorID(1)}

	msgs := []*qbft.SignedMessage{
		msg,
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change non unique signer",
		Pre:            pre,
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "invalid signed message: invalid signed message: non unique signer",
	}
}
