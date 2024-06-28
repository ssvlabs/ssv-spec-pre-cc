package startinstance

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// InvalidValue tests a starting an instance for an invalid value (not passing value check)
func InvalidValue() tests.SpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance invalid value",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:         testingutils.TestingInvalidValueCheck,
				ControllerPostRoot: "baf3ccea443a6c639b76dccf2d9c4fb5e48318473797de9b55e4d8de48fccc6b",
			},
		},
		ExpectedError: "value invalid: invalid value",
	}
}
