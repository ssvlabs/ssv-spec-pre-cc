package share

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// NoPartialQuorumDuplicate tests msg with < unique f+1 signers (but f+1 signers including duplicates)
func NoPartialQuorumDuplicate() *ShareTest {
	ks := testingutils.Testing4SharesSet()
	share := testingutils.TestingShare(ks)

	msg := &qbft.SignedMessage{
		Message: qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: testingutils.TestingIdentifier,
			Root:       testingutils.TestingQBFTRootData,
		},
		Signers: []types.OperatorID{1, 1},
	}

	return &ShareTest{
		Name:                     "no partial quorum duplicate",
		Share:                    *share,
		Message:                  *msg,
		ExpectedHasPartialQuorum: false,
		ExpectedHasQuorum:        false,
		ExpectedFullCommittee:    false,
		ExpectedError:            "non unique signer",
	}
}
