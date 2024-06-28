package consensusdata

import "github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"

// AttestationValidation tests a valid consensus data with AttestationData
func AttestationValidation() *ConsensusDataTest {
	return &ConsensusDataTest{
		Name:          "attestation validation",
		ConsensusData: *testingutils.TestAttesterConsensusData,
	}
}
