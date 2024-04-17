package partialsigcontainer

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

func Duplicate() tests.SpecTest {

	// Create a test key set
	ks := testingutils.Testing4SharesSet()

	// Create PartialSignatureMessage for testing
	msg1 := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg2 := testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight)
	msg3 := testingutils.PostConsensusAttestationMsg(ks.Shares[2], 2, qbft.FirstHeight)
	msgs := []*types.PartialSignatureMessage{msg1.Message.Messages[0], msg2.Message.Messages[0], msg3.Message.Messages[0]}

	return &PartialSigContainerTest{
		Name:            "duplicate",
		Quorum:          ks.Threshold,
		ValidatorPubKey: ks.ValidatorPK.Serialize(),
		SignatureMsgs:   msgs,
		ExpectedError:   "could not reconstruct a valid signature",
		ExpectedQuorum:  false,
	}
}
