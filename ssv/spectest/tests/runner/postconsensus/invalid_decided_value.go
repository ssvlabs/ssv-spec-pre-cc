package postconsensus

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// InvalidDecidedValue tests an invalid decided value
func InvalidDecidedValue() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	consensusDataByts := func(role types.BeaconRole) []byte {
		cd := &types.ConsensusData{
			Duty: types.Duty{
				Type:                    100,
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

	expectedErr := "failed processing post consensus message: invalid post-consensus message: no decided value"
	return &tests.MultiMsgProcessingSpecTest{
		Name: "post consensus decided invalid value",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee contribution",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3))),

					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
							[]*bls.SecretKey{
								ks.Shares[1], ks.Shares[2], ks.Shares[3],
							},
							[]types.OperatorID{1, 2, 3},
							qbft.Height(testingutils.TestingDutySlot),
							testingutils.SyncCommitteeContributionMsgID,
							consensusDataByts(types.BNRoleSyncCommitteeContribution),
						), nil)),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks))), // no qbft msg to mock the missing decided value
				},
				PostDutyRunnerStateRoot: "aff4af0dbbead81d6cb9dd4ff734d4660712a8b5ab8e9016a3f86e2c2ead7549",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgSyncCommittee(
						testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
							[]*bls.SecretKey{
								ks.Shares[1], ks.Shares[2], ks.Shares[3],
							},
							[]types.OperatorID{1, 2, 3},
							qbft.Height(testingutils.TestingDutySlot),
							testingutils.SyncCommitteeMsgID,
							consensusDataByts(types.BNRoleSyncCommittee),
						), nil)),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgSyncCommittee(nil, testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1))),
				},
				PostDutyRunnerStateRoot: "90c84430996225da29d9ed64d038a81d754599ce67d2a46a92689f2d4d57dfde",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedErr,
			},
			{
				Name:   "aggregator",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[2], ks.Shares[2], 2, 2))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[3], ks.Shares[3], 3, 3))),

					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgAggregator(
						testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
							[]*bls.SecretKey{
								ks.Shares[1], ks.Shares[2], ks.Shares[3],
							},
							[]types.OperatorID{1, 2, 3},
							qbft.Height(testingutils.TestingDutySlot),
							testingutils.AggregatorMsgID,
							consensusDataByts(types.BNRoleAggregator),
						), nil)),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgAggregator(nil, testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1))),
				},
				PostDutyRunnerStateRoot: "a93047858c5597f2b1de078a566e6b0227a217e10758a741ff7a2ed9e0a87d96",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "proposer",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, spec.DataVersionDeneb))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], ks.Shares[2], 2, 2, spec.DataVersionDeneb))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], ks.Shares[3], 3, 3, spec.DataVersionDeneb))),

					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(
						testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
							[]*bls.SecretKey{
								ks.Shares[1], ks.Shares[2], ks.Shares[3],
							},
							[]types.OperatorID{1, 2, 3},
							qbft.Height(testingutils.TestingDutySlot),
							testingutils.ProposerMsgID,
							consensusDataByts(types.BNRoleProposer),
						), nil)),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
				},
				PostDutyRunnerStateRoot: "0c965c41a9318297ad03c27e79ca2d2d0fee357fff21995c014182ce5e2970b3",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(testingutils.Testing4SharesSet().Shares[1], 1, spec.DataVersionDeneb),
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "proposer (blinded block)",
				Runner: testingutils.ProposerBlindedBlockRunner(ks),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[1], ks.Shares[1], 1, 1, spec.DataVersionDeneb))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[2], ks.Shares[2], 2, 2, spec.DataVersionDeneb))),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsgV(ks.Shares[3], ks.Shares[3], 3, 3, spec.DataVersionDeneb))),

					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(
						testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
							[]*bls.SecretKey{
								ks.Shares[1], ks.Shares[2], ks.Shares[3],
							},
							[]types.OperatorID{1, 2, 3},
							qbft.Height(testingutils.TestingDutySlot),
							testingutils.ProposerMsgID,
							consensusDataByts(types.BNRoleProposer),
						), nil)),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
				},
				PostDutyRunnerStateRoot: "d03cef76867bcb4540191a8e93a735b460ce5844271f718508a3821c404331a2",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(testingutils.Testing4SharesSet().Shares[1], 1, spec.DataVersionDeneb),
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: []*types.SignedSSVMessage{
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgAttester(
						testingutils.TestingCommitMultiSignerMessageWithHeightIdentifierAndFullData(
							[]*bls.SecretKey{
								ks.Shares[1], ks.Shares[2], ks.Shares[3],
							},
							[]types.OperatorID{1, 2, 3},
							qbft.Height(testingutils.TestingDutySlot),
							testingutils.AttesterMsgID,
							consensusDataByts(types.BNRoleAttester),
						), nil)),
					testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgAttester(nil, testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1,
						testingutils.TestingDutySlot))),
				},
				PostDutyRunnerStateRoot: "33953714dd71325c2ad309b2e122bf5fab016a5a2f1bfbf91125b3866c9dc844",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				ExpectedError:           expectedErr,
			},
		},
	}
}
