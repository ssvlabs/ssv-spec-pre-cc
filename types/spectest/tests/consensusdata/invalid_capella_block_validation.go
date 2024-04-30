package consensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// InvalidCapellaBlockValidation tests an invalid consensus data with capella block
func InvalidCapellaBlockValidation() *ConsensusDataTest {

	version := spec.DataVersionCapella

	cd := &types.ConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(version),
		Version: version,
		DataSSZ: []byte{},
	}
	return &ConsensusDataTest{
		Name:          "invalid capella block",
		ConsensusData: *cd,
		ExpectedError: "could not unmarshal ssz: incorrect size",
	}
}
