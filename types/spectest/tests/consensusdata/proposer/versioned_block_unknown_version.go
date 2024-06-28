package consensusdataproposer

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// VersionedBlockUnknownVersion tests a valid consensus data with unknown block
func VersionedBlockUnknownVersion() *ProposerSpecTest {
	unknownDataVersion := spec.DataVersion(100)
	cd := &types.ConsensusData{
		Duty:    *testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
		Version: unknownDataVersion,
		DataSSZ: testingutils.TestProposerBlindedBlockConsensusDataBytsV(spec.DataVersionDeneb),
	}

	cdSSZ, err := cd.MarshalSSZ()
	if err != nil {
		panic(err.Error())
	}

	return &ProposerSpecTest{
		Name:          "consensus data versioned block unknown version",
		DataCd:        cdSSZ,
		ExpectedError: fmt.Sprintf("unknown block version %s", unknownDataVersion.String()),
	}
}
