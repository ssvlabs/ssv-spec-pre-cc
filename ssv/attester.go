package ssv

import (
	"crypto/sha256"
	"encoding/json"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/go-bitfield"

	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
)

type AttesterRunner struct {
	BaseRunner *BaseRunner

	beacon         BeaconNode
	network        Network
	signer         types.KeyManager
	operatorSigner types.OperatorSigner
	valCheck       qbft.ProposedValueCheckF
}

func NewAttesterRunner(
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
	return &AttesterRunner{
		BaseRunner: &BaseRunner{
			BeaconRoleType:     types.BNRoleAttester,
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

func (r *AttesterRunner) StartNewDuty(duty *types.Duty) error {
	return r.BaseRunner.baseStartNewDuty(r, duty)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *AttesterRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *AttesterRunner) ProcessPreConsensus(signedMsg *types.SignedPartialSignatureMessage) error {
	return errors.New("no pre consensus sigs required for attester role")
}

func (r *AttesterRunner) ProcessConsensus(signedMsg *qbft.SignedMessage) error {
	decided, decidedValue, err := r.BaseRunner.baseConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing consensus message")
	}

	// Decided returns true only once so if it is true it must be for the current running instance
	if !decided {
		return nil
	}

	attestationData, err := decidedValue.GetAttestationData()
	if err != nil {
		return errors.Wrap(err, "could not get attestation data")
	}

	// specific duty sig
	msg, err := r.BaseRunner.signBeaconObject(r, attestationData, decidedValue.Duty.Slot, types.DomainAttester)
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

func (r *AttesterRunner) ProcessPostConsensus(signedMsg *types.SignedPartialSignatureMessage) error {
	quorum, roots, err := r.BaseRunner.basePostConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing post consensus message")
	}

	if !quorum {
		return nil
	}

	attestationData, err := r.GetState().DecidedValue.GetAttestationData()
	if err != nil {
		return errors.Wrap(err, "could not get attestation data")
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

		duty := r.GetState().DecidedValue.Duty

		aggregationBitfield := bitfield.NewBitlist(r.GetState().DecidedValue.Duty.CommitteeLength)
		aggregationBitfield.SetBitAt(duty.ValidatorCommitteeIndex, true)
		signedAtt := &phase0.Attestation{
			Data:            attestationData,
			Signature:       specSig,
			AggregationBits: aggregationBitfield,
		}

		// broadcast
		if err := r.beacon.SubmitAttestation(signedAtt); err != nil {
			return errors.Wrap(err, "could not submit to Beacon chain reconstructed attestation")
		}
	}
	r.GetState().Finished = true
	return nil
}

func (r *AttesterRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	return []ssz.HashRoot{}, types.DomainError, errors.New("no expected pre consensus roots for attester")
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *AttesterRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	attestationData, err := r.GetState().DecidedValue.GetAttestationData()
	if err != nil {
		return nil, phase0.DomainType{}, errors.Wrap(err, "could not get attestation data")
	}

	return []ssz.HashRoot{attestationData}, types.DomainAttester, nil
}

// executeDuty steps:
// 1) get attestation data from BN
// 2) start consensus on duty + attestation data
// 3) Once consensus decides, sign partial attestation and broadcast
// 4) collect 2f+1 partial sigs, reconstruct and broadcast valid attestation sig to the BN
func (r *AttesterRunner) executeDuty(duty *types.Duty) error {
	// TODO - waitOneThirdOrValidBlock

	attData, ver, err := r.GetBeaconNode().GetAttestationData(duty.Slot, duty.CommitteeIndex)
	if err != nil {
		return errors.Wrap(err, "failed to get attestation data")
	}

	attDataByts, err := attData.MarshalSSZ()
	if err != nil {
		return errors.Wrap(err, "could not marshal attestation data")
	}

	input := &types.ConsensusData{
		Duty:    *duty,
		Version: ver,
		DataSSZ: attDataByts,
	}

	if err := r.BaseRunner.decide(r, input); err != nil {
		return errors.Wrap(err, "can't start new duty runner instance for duty")
	}
	return nil
}

func (r *AttesterRunner) GetBaseRunner() *BaseRunner {
	return r.BaseRunner
}

func (r *AttesterRunner) GetNetwork() Network {
	return r.network
}

func (r *AttesterRunner) GetBeaconNode() BeaconNode {
	return r.beacon
}

func (r *AttesterRunner) GetShare() *types.Share {
	return r.BaseRunner.Share
}

func (r *AttesterRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *AttesterRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return r.valCheck
}

func (r *AttesterRunner) GetSigner() types.KeyManager {
	return r.signer
}

// Encode returns the encoded struct in bytes or error
func (r *AttesterRunner) Encode() ([]byte, error) {
	return json.Marshal(r)
}

// Decode returns error if decoding failed
func (r *AttesterRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &r)
}

// GetRoot returns the root used for signing and verification
func (r *AttesterRunner) GetRoot() ([32]byte, error) {
	marshaledRoot, err := r.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode DutyRunnerState")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}
