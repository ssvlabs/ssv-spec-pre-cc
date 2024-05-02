package commit

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// MultiSignerNoOverlap tests a multi signer commit msg which doesn't overlap previous valid commits
func MultiSignerNoOverlap() tests.SpecTest {
	pre := testingutils.BaseInstance()
	ks := testingutils.Testing4SharesSet()

	msgs := []*qbft.SignedMessage{
		testingutils.TestingProposalMessage(ks.Shares[1], 1),

		testingutils.TestingPrepareMessage(ks.Shares[1], 1),
		testingutils.TestingPrepareMessage(ks.Shares[2], 2),
		testingutils.TestingPrepareMessage(ks.Shares[3], 3),

		testingutils.TestingCommitMessage(ks.Shares[1], 1),
		testingutils.TestingCommitMultiSignerMessage([]*bls.SecretKey{ks.Shares[2], ks.Shares[3]}, []types.OperatorID{2, 3}),
	}
	return &tests.MsgProcessingSpecTest{
		Name:          "multi signer, no overlap",
		Pre:           pre,
		PostRoot:      "4e569d9a6c0421d2bb69a4c544f8f1e67c73a129d4e6bd1304ddbae8812cfa38",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.TestingPrepareMessage(ks.Shares[1], 1),
			testingutils.TestingCommitMessage(ks.Shares[1], 1),
		},
		ExpectedError: "invalid signed message: msg allows 1 signer",
	}
}
