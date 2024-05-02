package consensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// InvalidDuty tests an invalid consensus data with invalid duty
func InvalidDuty() *ConsensusDataTest {

	cd := &types.ConsensusData{
		Duty: types.Duty{
			Type:   types.BeaconRole(100),
			PubKey: testingutils.TestingValidatorPubKey,
		},
		Version: spec.DataVersionCapella,
		DataSSZ: testingutils.TestingAttestationDataBytes,
	}

	return &ConsensusDataTest{
		Name:          "invalid duty",
		ConsensusData: *cd,
		ExpectedError: "unknown duty role",
	}
}
