package consensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// FutureMessage tests a valid proposal future msg
func FutureMessage() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	futureMsgF := func(obj *types.ConsensusData, id []byte) *qbft.SignedMessage {
		fullData, _ := obj.Encode()
		root, _ := qbft.HashDataRoot(fullData)
		msg := &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     10,
			Round:      qbft.FirstRound,
			Identifier: id,
			Root:       root,
		}
		signed := testingutils.SignQBFTMsg(ks.Shares[1], 1, msg)
		signed.FullData = fullData

		return signed
	}

	const expectedError = "failed processing consensus message: future msg from height, could not process"

	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus future message",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgSyncCommitteeContribution(
						futureMsgF(testingutils.TestContributionProofWithJustificationsConsensusData(ks), testingutils.SyncCommitteeContributionMsgID),
						nil)),
				},
				PostDutyRunnerStateRoot: "68fd25b1cb30902e7b7b3e7ff674c3862ff956954a06fac0df485961b8bb3934",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgSyncCommittee(
						futureMsgF(testingutils.TestSyncCommitteeConsensusData, testingutils.SyncCommitteeMsgID),
						nil)),
				},
				PostDutyRunnerStateRoot: "022d7b2bbb97b88684ae317f5c2aaa46a56d1d272a65ffcbeb935d0511bbe7e0",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgAggregator(
						futureMsgF(testingutils.TestSelectionProofWithJustificationsConsensusData(ks), testingutils.AggregatorMsgID),
						nil)),
				},
				PostDutyRunnerStateRoot: "bdc7c2150e0f2d4669e112848f5140b52aba0367b60ff2b594d5a5bef3587834",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(
						futureMsgF(testingutils.TestProposerWithJustificationsConsensusDataV(ks, spec.DataVersionDeneb), testingutils.ProposerMsgID),
						nil)),
				},
				PostDutyRunnerStateRoot: "32dd1d1d7a4c34bb7dafc0866f69eb49f6a0a23755b135f83ad14d12e39fff82",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name:   "proposer (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(
						futureMsgF(testingutils.TestProposerBlindedWithJustificationsConsensusDataV(ks, spec.DataVersionDeneb), testingutils.ProposerMsgID),
						nil)),
				},
				PostDutyRunnerStateRoot: "58b946451dc5ccbd52fbc9e6bbe0ac888253d1708be018a3ff0b07762dd28891",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgAttester(
						futureMsgF(testingutils.TestAttesterConsensusData, testingutils.AttesterMsgID),
						nil)),
				},
				PostDutyRunnerStateRoot: "8ccbad4587df73b4a94e4c5d1c47c7ebfbc8e4e949518443a56f0f11d3ab70cd",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           expectedError,
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgValidatorRegistration(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.Shares[1], types.OperatorID(1),
							testingutils.ValidatorRegistrationMsgID, testingutils.TestAttesterConsensusDataByts,
							qbft.Height(testingutils.TestingDutySlot)),
						nil)),
				},
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "no consensus phase for validator registration",
			},
			{
				Name:   "voluntary exit",
				Runner: testingutils.VoluntaryExitRunner(ks),
				Duty:   &testingutils.TestingVoluntaryExitDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgVoluntaryExit(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(ks.Shares[1], types.OperatorID(1),
							testingutils.VoluntaryExitMsgID, testingutils.TestAttesterConsensusDataByts,
							qbft.Height(testingutils.TestingDutySlot)),
						nil)),
				},
				PostDutyRunnerStateRoot: "2ac409163b617c79a2a11d3919d6834d24c5c32f06113237a12afcf43e7757a0",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusVoluntaryExitMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "no consensus phase for voluntary exit",
			},
		},
	}
}
