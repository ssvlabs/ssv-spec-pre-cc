package consensusdata

import "github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"

// InvalidAttestationValidation tests an invalid consensus data with AttestationData
func InvalidAttestationValidation() *ConsensusDataTest {

	cd := testingutils.TestAttesterConsensusData
	cd.DataSSZ = testingutils.TestAggregatorConsensusDataByts

	return &ConsensusDataTest{
		Name:          "invalid attestation",
		ConsensusData: *cd,
		ExpectedError: "could not unmarshal ssz: incorrect size",
	}
}
