package qbft

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
)

// uponProposal process proposal message
// Assumes proposal message is valid!
func (i *Instance) uponProposal(signedProposal *SignedMessage, proposeMsgContainer *MsgContainer) error {
	addedMsg, err := proposeMsgContainer.AddFirstMsgForSignerAndRound(signedProposal)
	if err != nil {
		return errors.Wrap(err, "could not add proposal msg to container")
	}
	if !addedMsg {
		return nil // uponProposal was already called
	}

	newRound := signedProposal.Message.Round
	i.State.ProposalAcceptedForCurrentRound = signedProposal

	// A future justified proposal should bump us into future round and reset timer
	if signedProposal.Message.Round > i.State.Round {
		i.config.GetTimer().TimeoutForRound(signedProposal.Message.Round)
	}
	i.State.Round = newRound

	// value root
	r, err := HashDataRoot(signedProposal.FullData)
	if err != nil {
		return errors.Wrap(err, "could not hash input data")
	}

	prepare, err := CreatePrepare(i.State, i.config, newRound, r)
	if err != nil {
		return errors.Wrap(err, "could not create prepare msg")
	}

	if err := i.Broadcast(prepare); err != nil {
		return errors.Wrap(err, "failed to broadcast prepare message")
	}

	return nil
}

func isValidProposal(
	state *State,
	config IConfig,
	signedProposal *SignedMessage,
	valCheck ProposedValueCheckF,
	operators []*types.Operator,
) error {
	if signedProposal.Message.MsgType != ProposalMsgType {
		return errors.New("msg type is not proposal")
	}
	if signedProposal.Message.Height != state.Height {
		return errors.New("wrong msg height")
	}
	if len(signedProposal.GetSigners()) != 1 {
		return errors.New("msg allows 1 signer")
	}

	if !signedProposal.CheckSignersInCommittee(state.Share.Committee) {
		return errors.New("signer not in committee")
	}

	if !signedProposal.MatchedSigners([]types.OperatorID{proposer(state, config, signedProposal.Message.Round)}) {
		return errors.New("proposal leader invalid")
	}

	if err := signedProposal.Validate(); err != nil {
		return errors.Wrap(err, "proposal invalid")
	}

	// verify full data integrity
	r, err := HashDataRoot(signedProposal.FullData)
	if err != nil {
		return errors.Wrap(err, "could not hash input data")
	}
	if !bytes.Equal(signedProposal.Message.Root[:], r[:]) {
		return errors.New("H(data) != root")
	}

	// get justifications
	roundChangeJustification, _ := signedProposal.Message.GetRoundChangeJustifications() // no need to check error, checked on signedProposal.Validate()
	prepareJustification, _ := signedProposal.Message.GetPrepareJustifications()         // no need to check error, checked on signedProposal.Validate()

	if err := isProposalJustification(
		state,
		config,
		roundChangeJustification,
		prepareJustification,
		state.Height,
		signedProposal.Message.Round,
		signedProposal.FullData,
		valCheck,
	); err != nil {
		return errors.Wrap(err, "proposal not justified")
	}

	if (state.ProposalAcceptedForCurrentRound == nil && signedProposal.Message.Round == state.Round) ||
		signedProposal.Message.Round > state.Round {
		return nil
	}
	return errors.New("proposal is not valid with current state")
}

// isProposalJustification returns nil if the proposal and round change messages are valid and justify a proposal message for the provided round, value and leader
func isProposalJustification(
	state *State,
	config IConfig,
	roundChangeMsgs []*SignedMessage,
	prepareMsgs []*SignedMessage,
	height Height,
	round Round,
	fullData []byte,
	valCheck ProposedValueCheckF,
) error {
	if err := valCheck(fullData); err != nil {
		return errors.Wrap(err, "proposal fullData invalid")
	}

	if round == FirstRound {
		return nil
	} else {
		// check all round changes are valid for height and round
		// no quorum, duplicate signers,  invalid still has quorum, invalid no quorum
		// prepared
		for _, rc := range roundChangeMsgs {
			if err := validRoundChangeForDataVerifySignature(state, config, rc, height, round, fullData); err != nil {
				return errors.Wrap(err, "change round msg not valid")
			}
		}

		// check there is a quorum
		if !HasQuorum(state.Share, roundChangeMsgs) {
			return errors.New("change round has no quorum")
		}

		// previouslyPreparedF returns true if any on the round change messages have a prepared round and fullData
		previouslyPrepared, err := func(rcMsgs []*SignedMessage) (bool, error) {
			for _, rc := range rcMsgs {
				if rc.Message.RoundChangePrepared() {
					return true, nil
				}
			}
			return false, nil
		}(roundChangeMsgs)
		if err != nil {
			return errors.Wrap(err, "could not calculate if previously prepared")
		}

		if !previouslyPrepared {
			return nil
		} else {

			// check prepare quorum
			if !HasQuorum(state.Share, prepareMsgs) {
				return errors.New("prepares has no quorum")
			}

			// get a round change data for which there is a justification for the highest previously prepared round
			rcm, err := highestPrepared(roundChangeMsgs)
			if err != nil {
				return errors.Wrap(err, "could not get highest prepared")
			}
			if rcm == nil {
				return errors.New("no highest prepared")
			}

			// proposed fullData must equal highest prepared fullData
			r, err := HashDataRoot(fullData)
			if err != nil {
				return errors.Wrap(err, "could not hash input data")
			}
			if !bytes.Equal(r[:], rcm.Message.Root[:]) {
				return errors.New("proposed data doesn't match highest prepared")
			}

			// validate each prepare message against the highest previously prepared fullData and round
			for _, pm := range prepareMsgs {
				if err := validSignedPrepareForHeightRoundAndRootVerifySignature(
					config,
					pm,
					height,
					rcm.Message.DataRound,
					rcm.Message.Root,
					state.Share.Committee,
				); err != nil {
					return errors.New("signed prepare not valid")
				}
			}
			return nil
		}
	}
}

func proposer(state *State, config IConfig, round Round) types.OperatorID {
	// TODO - https://github.com/ConsenSys/qbft-formal-spec-and-verification/blob/29ae5a44551466453a84d4d17b9e083ecf189d97/dafny/spec/L1/node_auxiliary_functions.dfy#L304-L323
	return config.GetProposerF()(state, round)
}

// CreateProposal
/**
  	Proposal(
                        signProposal(
                            UnsignedProposal(
                                |current.blockchain|,
                                newRound,
                                digest(block)),
                            current.id),
                        block,
                        extractSignedRoundChanges(roundChanges),
                        extractSignedPrepares(prepares));
*/
func CreateProposal(state *State, config IConfig, fullData []byte, roundChanges, prepares []*SignedMessage) (*SignedMessage, error) {
	r, err := HashDataRoot(fullData)
	if err != nil {
		return nil, errors.Wrap(err, "could not hash input data")
	}

	roundChangesData, err := MarshalJustifications(roundChanges)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal justifications")
	}
	preparesData, err := MarshalJustifications(prepares)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal justifications")
	}

	msg := &Message{
		MsgType:    ProposalMsgType,
		Height:     state.Height,
		Round:      state.Round,
		Identifier: state.ID,

		Root:                     r,
		RoundChangeJustification: roundChangesData,
		PrepareJustification:     preparesData,
	}
	sig, err := config.GetShareSigner().SignRoot(msg, types.QBFTSignatureType, state.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed signing prepare msg")
	}

	signedMsg := &SignedMessage{
		Signature: sig,
		Signers:   []types.OperatorID{state.Share.OperatorID},
		Message:   *msg,

		FullData: fullData,
	}
	return signedMsg, nil
}
