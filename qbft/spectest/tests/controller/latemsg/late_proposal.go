package latemsg

import (
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils/comparable"
)

// LateProposal tests process late proposal msg for an instance which just decided
func LateProposal() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()
	sc := lateProposalStateComparison()

	msgs := testingutils.DecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData,
		testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, qbft.FirstHeight, ks)
	msgs = append(msgs, testingutils.TestingProposalMessage(ks.Shares[1], 1))

	return &tests.ControllerSpecTest{
		Name: "late proposal",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:    []byte{1, 2, 3, 4},
				InputMessages: msgs,
				ExpectedDecidedState: tests.DecidedState{
					DecidedVal: testingutils.TestingQBFTFullData,
					DecidedCnt: 1,
					BroadcastedDecided: testingutils.TestingCommitMultiSignerMessage(
						[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
						[]types.OperatorID{1, 2, 3},
					),
				},
				ControllerPostRoot:  sc.Root(),
				ControllerPostState: sc.ExpectedState,
			},
		},
		ExpectedError: "could not process msg: invalid signed message: proposal is not valid with current state",
	}
}

// LateProposalStateComparison returns the expected state comparison for LateProposal test.
// The controller is initialized with 4 shares and deciding messages for the first height from 3 nodes.
// The late msg is invalid since it is a proposal for the first height and the instance already decided, so it is ignored.
func lateProposalStateComparison() *comparable.StateComparison {
	ks := testingutils.Testing4SharesSet()
	msgs := testingutils.ExpectedDecidingMsgsForHeightWithRoot(testingutils.TestingQBFTRootData, testingutils.TestingQBFTFullData, testingutils.TestingIdentifier, qbft.FirstHeight, ks)

	contr := testingutils.NewTestingQBFTController(
		testingutils.TestingIdentifier,
		testingutils.TestingShare(testingutils.Testing4SharesSet()),
		testingutils.TestingConfig(testingutils.Testing4SharesSet()),
	)

	instance := &qbft.Instance{
		StartValue: []byte{1, 2, 3, 4},
		State: &qbft.State{
			Share:                           testingutils.TestingShare(testingutils.Testing4SharesSet()),
			ID:                              testingutils.TestingIdentifier,
			ProposalAcceptedForCurrentRound: testingutils.TestingProposalMessage(ks.Shares[1], types.OperatorID(1)),
			LastPreparedRound:               qbft.FirstRound,
			LastPreparedValue:               testingutils.TestingQBFTFullData,
			Decided:                         true,
			DecidedValue:                    testingutils.TestingQBFTFullData,
			Round:                           qbft.FirstRound,
		},
	}
	comparable.SetSignedMessages(instance, msgs)
	contr.StoredInstances = append(contr.StoredInstances, instance)

	return &comparable.StateComparison{ExpectedState: contr}
}
