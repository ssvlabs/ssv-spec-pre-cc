package newduty

import (
	"crypto/sha256"

	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/ssv"
	"github.com/ssvlabs/ssv-spec-pre-cc/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// PostInvalidDecided tests starting a new duty after prev was decided with an invalid decided value
func PostInvalidDecided() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	consensusDataByts := func(role types.BeaconRole) []byte {
		cd := &types.ConsensusData{
			Duty: types.Duty{
				Type:                    100, // invalid
				PubKey:                  testingutils.TestingValidatorPubKey,
				Slot:                    testingutils.TestingDutySlot,
				ValidatorIndex:          testingutils.TestingValidatorIndex,
				CommitteeIndex:          3,
				CommitteesAtSlot:        36,
				CommitteeLength:         128,
				ValidatorCommitteeIndex: 11,
			},
			Version: spec.DataVersionPhase0,
		}
		byts, _ := cd.Encode()
		return byts
	}

	// https://github.com/ssvlabs/ssv-spec-pre-cc/issues/285. We initialize the runner with an impossible decided value.
	// Maybe we should ensure that `ValidateDecided()` doesn't let the runner enter this state and delete the test?
	decideWrong := func(r ssv.Runner, duty *types.Duty) ssv.Runner {
		r.GetBaseRunner().State = ssv.NewRunnerState(3, duty)
		r.GetBaseRunner().State.RunningInstance = qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.Height(duty.Slot))
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, r.GetBaseRunner().State.RunningInstance)
		r.GetBaseRunner().QBFTController.Height = qbft.Height(duty.Slot)

		r.GetBaseRunner().State.RunningInstance.State.Decided = true
		decidedValue := sha256.Sum256(consensusDataByts(r.GetBaseRunner().BeaconRoleType))
		r.GetBaseRunner().State.RunningInstance.State.DecidedValue = decidedValue[:]

		return r
	}

	return &MultiStartNewRunnerDutySpecTest{
		Name: "new duty post invalid decided",
		Tests: []*StartNewRunnerDutySpecTest{
			{
				Name: "sync committee aggregator",
				Runner: decideWrong(testingutils.SyncCommitteeContributionRunner(ks),
					&testingutils.TestingSyncCommitteeContributionDuty),
				Duty:                    &testingutils.TestingSyncCommitteeContributionNexEpochDuty,
				PostDutyRunnerStateRoot: "4112802181d740f78b68b0c67e4220e689af2ac1011a51d0f4c10c4df315fac5",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  decideWrong(testingutils.SyncCommitteeRunner(ks), &testingutils.TestingSyncCommitteeDuty),
				Duty:                    &testingutils.TestingSyncCommitteeDutyNextEpoch,
				PostDutyRunnerStateRoot: "b3a2925d737a16363053e430c3ce317232fdff0a951e51bdf37cf593bfee0ae7",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:                    "aggregator",
				Runner:                  decideWrong(testingutils.AggregatorRunner(ks), &testingutils.TestingAggregatorDuty),
				Duty:                    &testingutils.TestingAggregatorDutyNextEpoch,
				PostDutyRunnerStateRoot: "b0ec12e65623dd1a95203d4a0a753bb6758f6bed3141a467bb7955530ae35ded",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofNextEpochMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:                    "proposer",
				Runner:                  decideWrong(testingutils.ProposerRunner(ks), testingutils.TestingProposerDutyV(spec.DataVersionDeneb)),
				Duty:                    testingutils.TestingProposerDutyNextEpochV(spec.DataVersionDeneb),
				PostDutyRunnerStateRoot: "c002484c2c25f5d97f625b5923484a062bdadb4eb21be9715dd9ae454883d890",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoNextEpochMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
					// broadcasts when starting a new duty
				},
			},
			{
				Name:                    "attester",
				Runner:                  decideWrong(testingutils.AttesterRunner(ks), &testingutils.TestingAttesterDuty),
				Duty:                    &testingutils.TestingAttesterDutyNextEpoch,
				PostDutyRunnerStateRoot: "679d7b60bd8bd84697331c6c872658a0cc8e3371156c7511be403d42abe18620",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
		},
	}
}
