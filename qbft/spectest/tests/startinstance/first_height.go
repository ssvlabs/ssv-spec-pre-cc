package startinstance

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// FirstHeight tests a starting the first instance
func FirstHeight() tests.SpecTest {
	return &tests.ControllerSpecTest{
		Name: "start instance first height",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: nil,
				},
				ControllerPostRoot: "47713c38fe74ce55959980781287886c603c2117a14dc8abce24dcb9be0093af",
				ExpectedTimerState: &testingutils.TimerState{
					Timeouts: 1,
					Round:    qbft.FirstRound,
				},
			},
		},
	}
}
