package consensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/ssv"
	ssvcomparable "github.com/ssvlabs/ssv-spec-pre-cc/ssv/spectest/comparable"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils/comparable"
)

// postDecidedSyncCommitteeContributionSC returns runner with decided instance and an extra commit message.
// There are pre-consensus messages in the container but no post-consensus messages.
func postDecidedSyncCommitteeContributionSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestSyncCommitteeContributionConsensusData
	cdBytes := testingutils.TestSyncCommitteeContributionConsensusDataByts

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.SyncCommitteeContributionRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					testingutils.SignedSSVMessageListF(ks, testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleSyncCommitteeContribution)[:3]),
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{},
				),
				DecidedValue: comparable.FixIssue178(cd, spec.DataVersionPhase0),
				StartingDuty: &cd.Duty,
				Finished:     false,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            testingutils.TestingDutySlot,
					LastPreparedRound: qbft.FirstRound,
					LastPreparedValue: cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.Shares[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes,
						qbft.Height(testingutils.TestingDutySlot)),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: comparable.NoErrorEncoding(comparable.FixIssue178(cd, spec.DataVersionBellatrix)),
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				append(
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleSyncCommitteeContribution)[3:10],
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.TestingCommitMessageWithHeightIdentifierAndFullData(ks.Shares[4], types.OperatorID(4), testingutils.TestingDutySlot, testingutils.SyncCommitteeContributionMsgID, cdBytes),
						nil,
					),
				),
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = testingutils.TestingDutySlot
			return ret
		}(),
	}
}

// postDecidedSyncCommitteeSC returns runner with decided instance and an extra commit message.
// There are no pre-consensus or post-consensus messages in the container.
func postDecidedSyncCommitteeSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestSyncCommitteeConsensusData
	cdBytes := testingutils.TestSyncCommitteeConsensusDataByts

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.SyncCommitteeRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					testingutils.SignedSSVMessageListF(ks, testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleSyncCommittee)[:3]),
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{},
				),
				DecidedValue: comparable.FixIssue178(cd, spec.DataVersionPhase0),
				StartingDuty: &cd.Duty,
				Finished:     false,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            testingutils.TestingDutySlot,
					LastPreparedRound: qbft.FirstRound,
					LastPreparedValue: cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.Shares[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes,
						qbft.Height(testingutils.TestingDutySlot)),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				append(
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleSyncCommittee)[:7],
					testingutils.SSVMsgSyncCommittee(
						testingutils.TestingCommitMessageWithHeightIdentifierAndFullData(ks.Shares[4], types.OperatorID(4), testingutils.TestingDutySlot, testingutils.SyncCommitteeMsgID, cdBytes),
						nil,
					),
				),
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = testingutils.TestingDutySlot
			return ret
		}(),
	}
}

// postDecidedAggregatorSC returns runner with decided instance and an extra commit message.
// There are pre-consensus messages in the container but no post-consensus messages.
func postDecidedAggregatorSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestAggregatorConsensusData
	cdBytes := testingutils.TestAggregatorConsensusDataByts

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.AggregatorRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					testingutils.SignedSSVMessageListF(ks, testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleAggregator)[:3]),
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{},
				),
				DecidedValue: comparable.FixIssue178(cd, spec.DataVersionPhase0),
				StartingDuty: &cd.Duty,
				Finished:     false,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            testingutils.TestingDutySlot,
					LastPreparedRound: qbft.FirstRound,
					LastPreparedValue: cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.Shares[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes,
						qbft.Height(testingutils.TestingDutySlot)),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				append(
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleAggregator)[3:10],
					testingutils.SSVMsgAggregator(
						testingutils.TestingCommitMessageWithHeightIdentifierAndFullData(ks.Shares[4], types.OperatorID(4), testingutils.TestingDutySlot, testingutils.AggregatorMsgID, cdBytes),
						nil,
					),
				),
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = testingutils.TestingDutySlot
			return ret
		}(),
	}
}

// postDecidedAttesterSC returns runner with decided instance and an extra commit message.
// There are no pre-consensus or post-consensus messages in the container.
func postDecidedAttesterSC() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestAttesterConsensusData
	cdBytes := testingutils.TestAttesterConsensusDataByts

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.AttesterRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					testingutils.SignedSSVMessageListF(ks, testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleAttester)[:3]),
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{},
				),
				DecidedValue: comparable.FixIssue178(cd, spec.DataVersionPhase0),
				StartingDuty: &cd.Duty,
				Finished:     false,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            testingutils.TestingDutySlot,
					LastPreparedRound: qbft.FirstRound,
					LastPreparedValue: cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.Shares[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes,
						qbft.Height(testingutils.TestingDutySlot)),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				append(
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleAttester)[:7],
					testingutils.SSVMsgAttester(
						testingutils.TestingCommitMessageWithHeightIdentifierAndFullData(ks.Shares[4], types.OperatorID(4), testingutils.TestingDutySlot, testingutils.AttesterMsgID, cdBytes),
						nil,
					),
				),
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = testingutils.TestingDutySlot
			return ret
		}(),
	}
}

// postDecidedProposerSC returns runner with decided instance and an extra commit message.
// There are pre-consensus messages in the container but no post-consensus messages.
func postDecidedProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestProposerConsensusDataV(version)
	cdBytes := testingutils.TestProposerConsensusDataBytsV(version)

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.ProposerRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					testingutils.SignedSSVMessageListF(ks, testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleProposer)[:3]),
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{},
				),
				DecidedValue: comparable.FixIssue178(cd, version),
				StartingDuty: &cd.Duty,
				Finished:     false,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.Height(testingutils.TestingDutySlotV(version)),
					LastPreparedRound: qbft.FirstRound,
					LastPreparedValue: cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.Shares[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes,
						qbft.Height(testingutils.TestingDutySlotV(version))),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				append(
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleProposer)[3:10],
					testingutils.SSVMsgProposer(
						testingutils.TestingCommitMessageWithHeightIdentifierAndFullData(ks.Shares[4], types.OperatorID(4),
							qbft.Height(testingutils.TestingDutySlotV(version)), testingutils.ProposerMsgID, cdBytes),
						nil,
					),
				),
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = qbft.Height(testingutils.TestingDutySlotV(version))
			return ret
		}(),
	}
}

// postDecidedBlindedProposerSC  returns runner with decided instance and an extra commit message.
// // There are pre-consensus messages in the container but no post-consensus messages.
func postDecidedBlindedProposerSC(version spec.DataVersion) *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	cd := testingutils.TestProposerBlindedBlockConsensusDataV(version)
	cdBytes := testingutils.TestProposerBlindedBlockConsensusDataBytsV(version)

	return &comparable.StateComparison{
		ExpectedState: func() ssv.Runner {
			ret := testingutils.ProposerBlindedBlockRunner(ks)
			ret.GetBaseRunner().State = &ssv.State{
				PreConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					testingutils.SignedSSVMessageListF(ks, testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleProposer)[:3]),
				),
				PostConsensusContainer: ssvcomparable.SetMessagesInContainer(
					ssv.NewPartialSigContainer(3),
					[]*types.SignedSSVMessage{},
				),
				DecidedValue: comparable.FixIssue178(cd, version),
				StartingDuty: &cd.Duty,
				Finished:     false,
			}
			ret.GetBaseRunner().State.RunningInstance = &qbft.Instance{
				State: &qbft.State{
					Share:             testingutils.TestingShare(ks),
					ID:                ret.GetBaseRunner().QBFTController.Identifier,
					Round:             qbft.FirstRound,
					Height:            qbft.Height(testingutils.TestingDutySlotV(version)),
					LastPreparedRound: qbft.FirstRound,
					LastPreparedValue: cdBytes,
					ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessageWithIdentifierAndFullData(
						ks.Shares[1], types.OperatorID(1), ret.GetBaseRunner().QBFTController.Identifier, cdBytes,
						qbft.Height(testingutils.TestingDutySlotV(version))),
					Decided:      true,
					DecidedValue: cdBytes,
				},
				StartValue: cdBytes,
			}
			comparable.SetMessages(
				ret.GetBaseRunner().State.RunningInstance,
				append(
					testingutils.ExpectedSSVDecidingMsgsV(cd, ks, types.BNRoleProposer)[3:10],
					testingutils.SSVMsgProposer(
						testingutils.TestingCommitMessageWithHeightIdentifierAndFullData(ks.Shares[4],
							types.OperatorID(4), qbft.Height(testingutils.TestingDutySlotV(version)), testingutils.ProposerMsgID,
							cdBytes),
						nil,
					),
				),
			)
			ret.GetBaseRunner().QBFTController.StoredInstances = append(ret.GetBaseRunner().QBFTController.StoredInstances, ret.GetBaseRunner().State.RunningInstance)
			ret.GetBaseRunner().QBFTController.Height = qbft.Height(testingutils.TestingDutySlotV(version))
			return ret
		}(),
	}
}
