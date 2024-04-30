package proposer

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"

	"github.com/ssvlabs/ssv-spec-pre-cc/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// ProposeRegularBlockDecidedBlinded tests proposing a regular block but the decided block is a blinded block. Full flow
func ProposeRegularBlockDecidedBlinded() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MsgProcessingSpecTest{
		Name:   "propose regular decide blinded",
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
					qbft.Height(testingutils.TestingDutySlotV(spec.DataVersionDeneb)),
					testingutils.ProposerMsgID,
					testingutils.TestProposerBlindedBlockConsensusDataBytsV(spec.DataVersionDeneb),
				), nil)),

			testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb))),
			testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[2], 2, spec.DataVersionDeneb))),
			testingutils.SignedSSVMessageF(ks, testingutils.SSVMsgProposer(nil, testingutils.PostConsensusProposerMsgV(ks.Shares[3], 3, spec.DataVersionDeneb))),
		},
		PostDutyRunnerStateRoot: "97b5fd3d786658f67d9d39df63a7a73b690e7873f9bdc107f6fcd401a42d98fc",
		OutputMessages: []*types.SignedPartialSignatureMessage{
			testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
			testingutils.PostConsensusProposerMsgV(ks.Shares[1], 1, spec.DataVersionDeneb),
		},
		BeaconBroadcastedRoots: []string{
			testingutils.GetSSZRootNoError(testingutils.TestingSignedBlindedBeaconBlockV(ks, spec.DataVersionDeneb)),
		},
	}
}
