package valcheckattestations

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/ssv/spectest/tests/valcheck"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// Valid tests valid data
func Valid() tests.SpecTest {
	return &valcheck.SpecTest{
		Name:       "attestation value check valid",
		Network:    types.PraterNetwork,
		BeaconRole: types.BNRoleAttester,
		Input:      testingutils.TestAttesterConsensusDataByts,
	}
}
