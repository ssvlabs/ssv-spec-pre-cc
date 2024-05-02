package consensusdata

import "github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"

// SyncCommitteeContributionValidation tests a valid consensus data with sync committee contrib.
func SyncCommitteeContributionValidation() *ConsensusDataTest {
	ks := testingutils.Testing4SharesSet()

	return &ConsensusDataTest{
		Name:          "sync committee contribution valid",
		ConsensusData: *testingutils.TestContributionProofWithJustificationsConsensusData(ks),
	}

}
