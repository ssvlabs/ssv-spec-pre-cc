package validation

import (
	"context"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/ssv"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
)

// MsgValidatorFunc represents a message validator
type MsgValidatorFunc = func(ctx context.Context, p peer.ID, msg *pubsub.Message) pubsub.ValidationResult

// To-Do: better spec the intended Message Validation behavior. Issue #375: https://github.com/ssvlabs/ssv-spec-pre-cc/issues/375
type MessageValidator interface {
	ValidateSignedSSVMessage(runner ssv.Runner, signedSSVMessage *types.SignedSSVMessage) error
	// Verifies the message's RSA signature with PKCSv1.5 encoding. Ref: https://datatracker.ietf.org/doc/html/rfc8017#section-8.2.2
	VerifySignature(runner ssv.Runner, signedSSVMessage *types.SignedSSVMessage) error
}

func MsgValidation(runner ssv.Runner, msgValidator MessageValidator) MsgValidatorFunc {
	return func(ctx context.Context, p peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
		signedSSVMsg, err := DecodePubsubMsg(msg)
		if err != nil {
			return pubsub.ValidationReject
		}

		// Message Validator
		if err := msgValidator.ValidateSignedSSVMessage(runner, signedSSVMsg); err != nil {
			return pubsub.ValidationReject
		}
		if err := msgValidator.VerifySignature(runner, signedSSVMsg); err != nil {
			return pubsub.ValidationReject
		}

		// Get SSVMessage
		ssvMsg, err := signedSSVMsg.GetSSVMessageFromData()
		if err != nil {
			return pubsub.ValidationReject
		}

		switch ssvMsg.GetType() {
		case types.SSVConsensusMsgType:
			if validateConsensusMsg(runner, ssvMsg.Data) != nil {
				return pubsub.ValidationReject
			}
		case types.SSVPartialSignatureMsgType:
			if validatePartialSigMsg(runner, ssvMsg.Data) != nil {
				return pubsub.ValidationReject
			}
		default:
			return pubsub.ValidationReject
		}

		return pubsub.ValidationAccept
	}
}

func DecodePubsubMsg(msg *pubsub.Message) (*types.SignedSSVMessage, error) {
	byts := msg.GetData()
	ret := &types.SignedSSVMessage{}
	if err := ret.Decode(byts); err != nil {
		return nil, err
	}
	return ret, nil
}

func validateConsensusMsg(runner ssv.Runner, data []byte) error {
	signedMsg := &qbft.SignedMessage{}
	if err := signedMsg.Decode(data); err != nil {
		return err
	}

	contr := runner.GetBaseRunner().QBFTController

	if err := contr.BaseMsgValidation(signedMsg); err != nil {
		return err
	}

	/**
	Main controller processing flow
	_______________________________
	All decided msgs are processed the same, out of instance
	All valid future msgs are saved in a container and can trigger highest decided futuremsg
	All other msgs (not future or decided) are processed normally by an existing instance (if found)
	*/
	if qbft.IsDecidedMsg(contr.Share, signedMsg) {
		return qbft.ValidateDecided(contr.GetConfig(), signedMsg, contr.Share)
	} else if signedMsg.Message.Height > contr.Height {
		return validateFutureMsg(contr.GetConfig(), signedMsg, contr.Share.Committee)
	} else {
		if inst := contr.StoredInstances.FindInstance(signedMsg.Message.Height); inst != nil {
			return inst.BaseMsgValidation(signedMsg)
		}
		return errors.New("unknown instance")
	}
}

func validatePartialSigMsg(runner ssv.Runner, data []byte) error {
	signedMsg := &types.SignedPartialSignatureMessage{}
	if err := signedMsg.Decode(data); err != nil {
		return err
	}

	if signedMsg.Message.Type == types.PostConsensusPartialSig {
		return runner.GetBaseRunner().ValidatePostConsensusMsg(runner, signedMsg)
	}
	return runner.GetBaseRunner().ValidatePreConsensusMsg(runner, signedMsg)
}

func validateFutureMsg(
	config qbft.IConfig,
	msg *qbft.SignedMessage,
	operators []*types.Operator,
) error {
	if err := msg.Validate(); err != nil {
		return errors.Wrap(err, "invalid decided msg")
	}

	if len(msg.GetSigners()) != 1 {
		return errors.New("allows 1 signer")
	}

	// verify signature
	if err := msg.Signature.VerifyByOperators(msg, config.GetSignatureDomainType(), types.QBFTSignatureType, operators); err != nil {
		return errors.Wrap(err, "msg signature invalid")
	}

	return nil
}
