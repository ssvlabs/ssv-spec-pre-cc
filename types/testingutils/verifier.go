package testingutils

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"

	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
)

type testingVerifier struct {
	signaturesCache map[types.OperatorID]map[[32]byte][256]byte
}

func NewTestingVerifier() types.SignatureVerifier {
	return &testingVerifier{
		signaturesCache: make(map[uint64]map[[32]byte][256]byte),
	}
}

func (v *testingVerifier) Verify(msg *types.SignedSSVMessage, operators []*types.Operator) error {

	// Get message hash
	hash := sha256.Sum256(msg.Data)

	// Find operator that matches ID with the signer and verify signature
	for _, op := range operators {
		// Find operator
		if op.OperatorID == msg.GetOperatorID() {

			// Check cache
			if v.HasSignature(op.OperatorID, hash, msg.Signature) {
				return nil
			}

			// Get public key
			parsedPk, err := x509.ParsePKIXPublicKey(op.SSVOperatorPubKey)
			if err != nil {
				return errors.Wrap(err, "could not parse signer public key")
			}
			pk, ok := parsedPk.(*rsa.PublicKey)
			if !ok {
				return errors.New("could not parse signer public key")
			}

			// Verify
			err = rsa.VerifyPKCS1v15(pk, crypto.SHA256, hash[:], msg.Signature[:])

			if err == nil {
				v.SaveSignature(op.OperatorID, hash, msg.Signature)
			}
			return err
		}
	}

	return errors.New("unknown signer")
}

func (v *testingVerifier) HasSignature(operatorID types.OperatorID, root [32]byte, signature [256]byte) bool {
	if _, found := v.signaturesCache[operatorID]; !found {
		return false
	}

	storedSignature, found := v.signaturesCache[operatorID][root]
	if !found {
		return false
	}

	return bytes.Equal(storedSignature[:], signature[:])
}

func (v *testingVerifier) SaveSignature(operatorID types.OperatorID, root [32]byte, signature [256]byte) {
	if _, found := v.signaturesCache[operatorID]; !found {
		v.signaturesCache[operatorID] = make(map[[32]byte][256]byte)
	}
	v.signaturesCache[operatorID][root] = signature
}

// Verifies a list of SignedSSVMessage using the operators list
func VerifyListOfSignedSSVMessages(msgs []*types.SignedSSVMessage, operators []*types.Operator) error {
	verifier := NewTestingVerifier()

	for _, msg := range msgs {
		err := verifier.Verify(msg, operators)
		if err != nil {
			return err
		}
	}
	return nil
}
