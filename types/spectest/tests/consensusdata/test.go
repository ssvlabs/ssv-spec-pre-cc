package consensusdata

import (
	comparable2 "github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils/comparable"
	reflect2 "reflect"
	"testing"

	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/stretchr/testify/require"
)

type ConsensusDataTest struct {
	Name          string
	ConsensusData types.ConsensusData
	ExpectedError string
}

func (test *ConsensusDataTest) TestName() string {
	return "consensusdata " + test.Name
}

func (test *ConsensusDataTest) Run(t *testing.T) {

	err := test.ConsensusData.Validate()

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, err, test.ExpectedError)
	} else {
		require.NoError(t, err)
	}

	comparable2.CompareWithJson(t, test, test.TestName(), reflect2.TypeOf(test).String())
}
