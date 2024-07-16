package consensusdata

import (
	"github.com/AKorpusenko/genesis-go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// DenebBlockValidation tests a valid consensus data with deneb block
func DenebBlockValidation() *ConsensusDataTest {
	ks := testingutils.Testing4SharesSet()

	return &ConsensusDataTest{
		Name:          "valid deneb block",
		ConsensusData: *testingutils.TestProposerWithJustificationsConsensusDataV(ks, spec.DataVersionDeneb),
	}
}
