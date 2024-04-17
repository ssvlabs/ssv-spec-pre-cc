package consensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// SyncCommitteeNoJustifications tests a valid consensus data with no sync committee pre-consensus justifications
func SyncCommitteeNoJustifications() *ConsensusDataTest {

	cd := types.ConsensusData{
		Duty:    testingutils.TestingSyncCommitteeDuty,
		Version: spec.DataVersionCapella,
		DataSSZ: testingutils.TestingSyncCommitteeBlockRoot[:],
	}

	return &ConsensusDataTest{
		Name:          "sync committee no pre-consensus justification",
		ConsensusData: cd,
	}
}
