package consensusdata

import "github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"

// InvalidSyncCommitteeContributionValidation tests an invalid consensus data with sync committee contrib.
func InvalidSyncCommitteeContributionValidation() *ConsensusDataTest {

	cd := testingutils.TestSyncCommitteeContributionConsensusData
	cd.DataSSZ = testingutils.TestingAttestationDataBytes

	return &ConsensusDataTest{
		Name:          "invalid sync committee contribution",
		ConsensusData: *cd,
		ExpectedError: "could not unmarshal ssz: four",
	}
}
