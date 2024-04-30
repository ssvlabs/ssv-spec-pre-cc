package ssv

import (
	"crypto/sha256"
	"encoding/json"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"

	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
)

type SyncCommitteeRunner struct {
	BaseRunner *BaseRunner

	beacon         BeaconNode
	network        Network
	signer         types.KeyManager
	operatorSigner types.OperatorSigner
	valCheck       qbft.ProposedValueCheckF
}

func NewSyncCommitteeRunner(
	beaconNetwork types.BeaconNetwork,
	share *types.Share,
	qbftController *qbft.Controller,
	beacon BeaconNode,
	network Network,
	signer types.KeyManager,
	operatorSigner types.OperatorSigner,
	valCheck qbft.ProposedValueCheckF,
	highestDecidedSlot phase0.Slot,
) Runner {
	return &SyncCommitteeRunner{
		BaseRunner: &BaseRunner{
			BeaconRoleType:     types.BNRoleSyncCommittee,
			BeaconNetwork:      beaconNetwork,
			Share:              share,
			QBFTController:     qbftController,
			highestDecidedSlot: highestDecidedSlot,
		},

		beacon:         beacon,
		network:        network,
		signer:         signer,
		operatorSigner: operatorSigner,
		valCheck:       valCheck,
	}
}

func (r *SyncCommitteeRunner) StartNewDuty(duty *types.Duty) error {
	return r.BaseRunner.baseStartNewDuty(r, duty)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *SyncCommitteeRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *SyncCommitteeRunner) ProcessPreConsensus(signedMsg *types.SignedPartialSignatureMessage) error {
	return errors.New("no pre consensus sigs required for sync committee role")
}

func (r *SyncCommitteeRunner) ProcessConsensus(signedMsg *qbft.SignedMessage) error {
	decided, decidedValue, err := r.BaseRunner.baseConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing consensus message")
	}

	// Decided returns true only once so if it is true it must be for the current running instance
	if !decided {
		return nil
	}

	// specific duty sig
	root, err := decidedValue.GetSyncCommitteeBlockRoot()
	if err != nil {
		return errors.Wrap(err, "could not get sync committee block root")
	}
	msg, err := r.BaseRunner.signBeaconObject(r, types.SSZBytes(root[:]), decidedValue.Duty.Slot, types.DomainSyncCommittee)
	if err != nil {
		return errors.Wrap(err, "failed signing attestation data")
	}
	postConsensusMsg := &types.PartialSignatureMessages{
		Type:     types.PostConsensusPartialSig,
		Slot:     decidedValue.Duty.Slot,
		Messages: []*types.PartialSignatureMessage{msg},
	}

	postSignedMsg, err := r.BaseRunner.signPostConsensusMsg(r, postConsensusMsg)
	if err != nil {
		return errors.Wrap(err, "could not sign post consensus msg")
	}

	data, err := postSignedMsg.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode post consensus signature msg")
	}

	ssvMsg := &types.SSVMessage{
		MsgType: types.SSVPartialSignatureMsgType,
		MsgID:   types.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey, r.BaseRunner.BeaconRoleType),
		Data:    data,
	}

	msgToBroadcast, err := types.SSVMessageToSignedSSVMessage(ssvMsg, r.BaseRunner.Share.OperatorID, r.operatorSigner.SignSSVMessage)
	if err != nil {
		return errors.Wrap(err, "could not create SignedSSVMessage from SSVMessage")
	}

	if err := r.GetNetwork().Broadcast(ssvMsg.GetID(), msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial post consensus sig")
	}
	return nil
}

func (r *SyncCommitteeRunner) ProcessPostConsensus(signedMsg *types.SignedPartialSignatureMessage) error {
	quorum, roots, err := r.BaseRunner.basePostConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing post consensus message")
	}

	if !quorum {
		return nil
	}

	blockRoot, err := r.GetState().DecidedValue.GetSyncCommitteeBlockRoot()
	if err != nil {
		return errors.Wrap(err, "could not get sync committee block root")
	}

	for _, root := range roots {
		sig, err := r.GetState().ReconstructBeaconSig(r.GetState().PostConsensusContainer, root, r.GetShare().ValidatorPubKey)
		if err != nil {
			// If the reconstructed signature verification failed, fall back to verifying each partial signature
			for _, root := range roots {
				r.BaseRunner.FallBackAndVerifyEachSignature(r.GetState().PostConsensusContainer, root)
			}
			return errors.Wrap(err, "got post-consensus quorum but it has invalid signatures")
		}
		specSig := phase0.BLSSignature{}
		copy(specSig[:], sig)

		msg := &altair.SyncCommitteeMessage{
			Slot:            r.GetState().DecidedValue.Duty.Slot,
			BeaconBlockRoot: blockRoot,
			ValidatorIndex:  r.GetState().DecidedValue.Duty.ValidatorIndex,
			Signature:       specSig,
		}
		if err := r.GetBeaconNode().SubmitSyncMessage(msg); err != nil {
			return errors.Wrap(err, "could not submit to Beacon chain reconstructed signed sync committee")
		}
	}
	r.GetState().Finished = true
	return nil
}

func (r *SyncCommitteeRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	return []ssz.HashRoot{}, types.DomainError, errors.New("no expected pre consensus roots for sync committee")
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *SyncCommitteeRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	root, err := r.GetState().DecidedValue.GetSyncCommitteeBlockRoot()
	if err != nil {
		return nil, phase0.DomainType{}, errors.Wrap(err, "could not get sync committee block root")
	}

	return []ssz.HashRoot{types.SSZBytes(root[:])}, types.DomainSyncCommittee, nil
}

// executeDuty steps:
// 1) get sync block root from BN
// 2) start consensus on duty + block root data
// 3) Once consensus decides, sign partial block root and broadcast
// 4) collect 2f+1 partial sigs, reconstruct and broadcast valid sync committee sig to the BN
func (r *SyncCommitteeRunner) executeDuty(duty *types.Duty) error {
	// TODO - waitOneThirdOrValidBlock

	root, ver, err := r.GetBeaconNode().GetSyncMessageBlockRoot(duty.Slot)
	if err != nil {
		return errors.Wrap(err, "failed to get sync committee block root")
	}

	input := &types.ConsensusData{
		Duty:    *duty,
		Version: ver,
		DataSSZ: root[:],
	}

	if err := r.BaseRunner.decide(r, input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}
	return nil
}

func (r *SyncCommitteeRunner) GetBaseRunner() *BaseRunner {
	return r.BaseRunner
}

func (r *SyncCommitteeRunner) GetNetwork() Network {
	return r.network
}

func (r *SyncCommitteeRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

func (r *SyncCommitteeRunner) GetShare() *types.Share {
	return r.BaseRunner.Share
}

func (r *SyncCommitteeRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *SyncCommitteeRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return r.valCheck
}

func (r *SyncCommitteeRunner) GetSigner() types.KeyManager {
	return r.signer
}

// Encode returns the encoded struct in bytes or error
func (r *SyncCommitteeRunner) Encode() ([]byte, error) {
	return json.Marshal(r)
}

// Decode returns error if decoding failed
func (r *SyncCommitteeRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &r)
}

// GetRoot returns the root used for signing and verification
func (r *SyncCommitteeRunner) GetRoot() ([32]byte, error) {
	marshaledRoot, err := r.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode DutyRunnerState")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}
