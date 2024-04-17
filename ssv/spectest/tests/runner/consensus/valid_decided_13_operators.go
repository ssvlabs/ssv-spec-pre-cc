package consensus

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// ValidDecided13Operators tests a valid decided value (13 operators)
func ValidDecided13Operators() tests.SpecTest {
	ks := testingutils.Testing13SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus valid decided 13 operators",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:                    "sync committee contribution",
				Runner:                  testingutils.SyncCommitteeContributionRunner(ks),
				Duty:                    &testingutils.TestingSyncCommitteeContributionDuty,
				Messages:                testingutils.SignedSSVMessageListF(ks, testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeContributionConsensusData, ks, types.BNRoleSyncCommitteeContribution)),
				PostDutyRunnerStateRoot: "2feec534f4ef14ea96e83b17d4feaf3f40e4739a72fc7118c7da1ef2c1331937",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusSyncCommitteeContributionMsg(ks.Shares[1], 1, ks),
				},
			},
			{
				Name:                    "sync committee",
				Runner:                  testingutils.SyncCommitteeRunner(ks),
				Duty:                    &testingutils.TestingSyncCommitteeDuty,
				Messages:                testingutils.SignedSSVMessageListF(ks, testingutils.SSVDecidingMsgsV(testingutils.TestSyncCommitteeConsensusData, ks, types.BNRoleSyncCommittee)),
				PostDutyRunnerStateRoot: "72c411326ede815faebad64a3d1d0ba4575d1e103f83d2008ac9a26327e4691c",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusSyncCommitteeMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "aggregator",
				Runner:                  testingutils.AggregatorRunner(ks),
				Duty:                    &testingutils.TestingAggregatorDuty,
				Messages:                testingutils.SignedSSVMessageListF(ks, testingutils.SSVDecidingMsgsV(testingutils.TestAggregatorConsensusData, ks, types.BNRoleAggregator)),
				PostDutyRunnerStateRoot: "dd79a9da2d025a61d4e4ebdebbc9a331a876a43bb8c8e7c274132ae1c4a35175",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1),
					testingutils.PostConsensusAggregatorMsg(ks.Shares[1], 1),
				},
			},
			{
				Name:                    "proposer",
				Runner:                  testingutils.ProposerRunner(ks),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages:                testingutils.SignedSSVMessageListF(ks, testingutils.SSVDecidingMsgsV(testingutils.TestProposerConsensusDataV(spec.DataVersionDeneb), ks, types.BNRoleProposer)),
				PostDutyRunnerStateRoot: "942d09a2aaaddf2ef7b2388fd1a9ffb94d590c645fc90c2c6b27fa5414afd61e",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
					testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
				},
			},
			{
				Name:                    "proposer (blinded block)",
				Runner:                  testingutils.ProposerBlindedBlockRunner(ks),
				Duty:                    testingutils.TestingProposerDutyV(spec.DataVersionDeneb),
				Messages:                testingutils.SignedSSVMessageListF(ks, testingutils.SSVDecidingMsgsV(testingutils.TestProposerBlindedBlockConsensusDataV(spec.DataVersionDeneb), ks, types.BNRoleProposer)),
				PostDutyRunnerStateRoot: "abf4eaf532eda91fd6db48136938985a655174ab5101e698d71737c0ed48c8c3",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
					testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
				},
			},
			{
				Name:                    "attester",
				Runner:                  testingutils.AttesterRunner(ks),
				Duty:                    &testingutils.TestingAttesterDuty,
				Messages:                testingutils.SignedSSVMessageListF(ks, testingutils.SSVDecidingMsgsV(testingutils.TestAttesterConsensusData, ks, types.BNRoleAttester)),
				PostDutyRunnerStateRoot: "7dc2522866171588a2f00aa40690832ee06b9f4f225cbfc97fbbb991f7cbc4a6",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PostConsensusAttestationMsg(ks.Shares[1], 1, qbft.FirstHeight),
				},
			},
		},
	}
}
