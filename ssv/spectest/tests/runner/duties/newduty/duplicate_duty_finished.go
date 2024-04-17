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

// DuplicateDutyFinished is a test that runs the following scenario:
//   - Runner is assigned a duty
//   - Runner finishes the duty
//   - Runner is assigned the same duty again
func DuplicateDutyFinished() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	finishRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.Height(duty.Slot))
		r.GetBaseRunner().State.RunningInstance.State.Decided = true
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)
		r.GetBaseRunner().QBFTController.Height = qbft.Height(duty.Slot)
		r.GetBaseRunner().State.Finished = true
		return r
	}

	expectedError := fmt.Sprintf("can't start duty: duty for slot %d already passed. Current height is %d",
		testingutils.TestingDutySlot,
		testingutils.TestingDutySlot)

	// finishTaskRunner is a helper function that finishes a task runner and returns it
	// task is an operation that isn't a beacon duty, e.g. validator registration
	finishTaskRunner := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.Finished = true
		return r
	}

	expectedTaskError := fmt.Sprintf("can't start non-beacon duty: duty for slot %d already passed. "+
		"Current slot is %d",
		testingutils.TestingDutySlot,
		testingutils.TestingDutySlot)

	return &MultiStartNewRunnerDutySpecTest{
		Name: "duplicate duty finished",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name:                    "sync committee aggregator",
				Runner:                  finishRunner(testingutils.SyncCommitteeContributionRunner(ks), &testingutils.TestingSyncCommitteeContributionDuty),
				Duty:                    &testingutils.TestingSyncCommitteeContributionDuty,
				PostDutyRunnerStateRoot: "c8ce3cec33a9e557f52c1392f96b613ed2d37b24b54a1c9429a7dbff91f212eb",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "sync committee",
				Runner:                  finishRunner(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:                    &testingutils.TestingSyncCommitteeDuty,
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				PostDutyRunnerStateRoot: "9ed70d234980d27628811e78f59b0a723ae2bd768ad9ce02943aa3fdf737e2c5",
				ExpectedError:           expectedError,
			},
			{
				Name:                    "aggregator",
				Runner:                  finishRunner(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty),
				Duty:                    &testingutils.TestingAggregatorDuty,
				PostDutyRunnerStateRoot: "3674c8986f519e022f76377d00c5d27ef2e53faaf6bffce4eb692bf5d387d6b2",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedError,
			},
			{
				Name:                    "proposer",
				Runner:                  finishRunner(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyV(spec.DataVersionDeneb)),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				PostDutyRunnerStateRoot: "a91e014950037e5dc2ab9e801d0170b90b82f592029a2409c2332f252368d71d",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionDeneb), // broadcasts when starting a new duty
				},
				ExpectedError: fmt.Sprintf("can't start duty: duty for slot %d already passed. Current height is %d",
					testingutils.TestingDutySlotV(spec.DataVersionDeneb),
					testingutils.TestingDutySlotV(spec.DataVersionDeneb)),
			},
			{
				Name:                    "attester",
				Runner:                  finishRunner(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:                    &testingutils.TestingAttesterDuty,
				PostDutyRunnerStateRoot: "a96148ae850dd3d3a0d63869a95702174739151fa271ba463a3c163cabe35e13",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedError,
			},
			{
				Name: "validator registration",
				Runner: finishTaskRunner(testingutils.ValidatorRegistrationRunner(ks),
					&testingutils.TestingValidatorRegistrationDuty),
				Duty:                    &testingutils.TestingValidatorRegistrationDuty,
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedTaskError,
			},
			{
				Name: "voluntary exit",
				Runner: finishTaskRunner(testingutils.VoluntaryExitRunner(ks),
					&testingutils.TestingVoluntaryExitDuty),
				Duty:                    &testingutils.TestingVoluntaryExitDuty,
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedTaskError,
			},
		},
	}
}
