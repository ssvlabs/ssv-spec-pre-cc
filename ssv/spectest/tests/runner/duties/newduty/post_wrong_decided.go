package newduty

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/ssv"
	"github.com/ssvlabs/ssv-spec-pre-cc/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// PostWrongDecided tests starting a new duty after prev was decided wrongly (future decided)
// This can happen if we receive a future decided message from the network.
func PostWrongDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	// https://github.com/ssvlabs/ssv-spec-pre-cc/issues/285. We initialize the runner with an impossible decided value.
	// Maybe we should ensure that `ValidateDecided()` doesn't let the runner enter this state and delete the test?
	decideWrong := func(r ssv.Runner, duty *types.Duty, higherDecidedSlot qbft.Height) ssv.Runner {
		storedInstances := r.GetBaseRunner().QBFTController.StoredInstances
		storedInstances = append(storedInstances, nil)
		storedInstances = append(storedInstances, nil)

		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight)
		r.GetBaseRunner().State.RunningInstance.State.Decided = true
		storedInstances[1] = r.GetBaseRunner().State.RunningInstance

		higherDecided := qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			higherDecidedSlot)
		higherDecided.State.Decided = true
		higherDecided.State.DecidedValue = []byte{1, 2, 3, 4}
		storedInstances[0] = higherDecided
		r.GetBaseRunner().QBFTController.Height = higherDecidedSlot
		// TODO: hacky fix to a bug in the test.
		// You can't append a copied slice and expect the original to change in go. Since maybe we want to delete
		// the test I didn't do it nicer.
		r.GetBaseRunner().QBFTController.StoredInstances = storedInstances
		return r
	}

	expectedError := fmt.Sprintf("can't start duty: duty for slot %d already passed. Current height is %d",
		testingutils.TestingDutySlot, 50)

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty post wrong decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  decideWrong(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionDuty, 50),
				Duty:                    &testingutils.TestingSyncCommitteeContributionDuty,
				PostDutyRunnerStateRoot: "4fce8afe24a8f812c9daccfb54c8247771c88b48d161b06901669d1e23ce7a0d",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "sync committee",
				Runner:                  decideWrong(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty, 50),
				Duty:                    &testingutils.TestingSyncCommitteeDuty,
				PostDutyRunnerStateRoot: "daf5d2d2469f71f2f9c08a75fb45f1be7f6cdc8b367dedf5adc50c066f73f942",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedError,
			},
			{
				Name:                    "aggregator",
				Runner:                  decideWrong(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty, 50),
				Duty:                    &testingutils.TestingAggregatorDuty,
				PostDutyRunnerStateRoot: "533fa28f89164dc4a26bdd4e2cd55cca8c17375d3e88251f2480ae727ced1ee1",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "proposer",
				Runner:                  decideWrong(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyV(spec.DataVersionDeneb), qbft.Height(testingutils.TestingDutySlotV(spec.DataVersionDeneb)+50)),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				PostDutyRunnerStateRoot: "fd8de7873c9cf83ec9366f856f4b8d81d77c215f935185d90cea8be8dcb44089",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
				ExpectedError: fmt.Sprintf("can't start duty: duty for slot %d already passed. Current height is %d",
					testingutils.TestingDutySlotV(spec.DataVersionDeneb), testingutils.TestingDutySlotV(spec.DataVersionDeneb)+50),
			},
			{
				Name:                    "attester",
				Runner:                  decideWrong(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty, 50),
				Duty:                    &testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "8efe8f69adeecb1da0762d54aeb7b2970eb293395ad6e379cf2521125881f5fe",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedError,
			},
		},
	}
}
