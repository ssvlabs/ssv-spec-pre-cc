package consensusdata

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// ProposerNoJustifications tests an invalid consensus data with no proposer justifications
func ProposerNoJustifications() *ConsensusDataTest {

	// To-do: add error when pre-consensus justification check is added.

	cd := testingutils.TestProposerConsensusDataV(spec.DataVersionCapella)
	cd.PreConsensusJustifications = nil

	return &ConsensusDataTest{
		Name:          "proposer no justification",
		ConsensusData: *cd,
	}
}
