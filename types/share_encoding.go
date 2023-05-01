// Code generated by fastssz. DO NOT EDIT.
// Hash: 9e7e6744320cd7be5a30a62f1b743fb2ab72ecc2d7ab2093e5ed4e7c959d767a
// Version: 0.1.2
package types

import (
	ssz "github.com/ferranbt/fastssz"
)

// MarshalSSZ ssz marshals the Share object
func (s *Share) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(s)
}

// MarshalSSZTo ssz marshals the Share object to a target array
func (s *Share) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(164)

	// Field (0) 'OperatorID'
	dst = ssz.MarshalUint64(dst, uint64(s.OperatorID))

	// Field (1) 'ValidatorPubKey'
	if size := len(s.ValidatorPubKey); size != 48 {
		err = ssz.ErrBytesLengthFn("Share.ValidatorPubKey", size, 48)
		return
	}
	dst = append(dst, s.ValidatorPubKey...)

	// Field (2) 'SharePubKey'
	if size := len(s.SharePubKey); size != 48 {
		err = ssz.ErrBytesLengthFn("Share.SharePubKey", size, 48)
		return
	}
	dst = append(dst, s.SharePubKey...)

	// Offset (3) 'Committee'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(s.Committee) * 56

	// Field (4) 'DomainType'
	dst = append(dst, s.DomainType[:]...)

	// Field (5) 'FeeRecipientAddress'
	dst = append(dst, s.FeeRecipientAddress[:]...)

	// Field (6) 'Graffiti'
	if size := len(s.Graffiti); size != 32 {
		err = ssz.ErrBytesLengthFn("Share.Graffiti", size, 32)
		return
	}
	dst = append(dst, s.Graffiti...)

	// Field (3) 'Committee'
	if size := len(s.Committee); size > 13 {
		err = ssz.ErrListTooBigFn("Share.Committee", size, 13)
		return
	}
	for ii := 0; ii < len(s.Committee); ii++ {
		if dst, err = s.Committee[ii].MarshalSSZTo(dst); err != nil {
			return
		}
	}

	return
}

// UnmarshalSSZ ssz unmarshals the Share object
func (s *Share) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 164 {
		return ssz.ErrSize
	}

	tail := buf
	var o3 uint64

	// Field (0) 'OperatorID'
	s.OperatorID = OperatorID(ssz.UnmarshallUint64(buf[0:8]))

	// Field (1) 'ValidatorPubKey'
	if cap(s.ValidatorPubKey) == 0 {
		s.ValidatorPubKey = make([]byte, 0, len(buf[8:56]))
	}
	s.ValidatorPubKey = append(s.ValidatorPubKey, buf[8:56]...)

	// Field (2) 'SharePubKey'
	if cap(s.SharePubKey) == 0 {
		s.SharePubKey = make([]byte, 0, len(buf[56:104]))
	}
	s.SharePubKey = append(s.SharePubKey, buf[56:104]...)

	// Offset (3) 'Committee'
	if o3 = ssz.ReadOffset(buf[104:108]); o3 > size {
		return ssz.ErrOffset
	}

	if o3 < 164 {
		return ssz.ErrInvalidVariableOffset
	}

	// Field (4) 'DomainType'
	copy(s.DomainType[:], buf[108:112])

	// Field (5) 'FeeRecipientAddress'
	copy(s.FeeRecipientAddress[:], buf[112:132])

	// Field (6) 'Graffiti'
	if cap(s.Graffiti) == 0 {
		s.Graffiti = make([]byte, 0, len(buf[132:164]))
	}
	s.Graffiti = append(s.Graffiti, buf[132:164]...)

	// Field (3) 'Committee'
	{
		buf = tail[o3:]
		num, err := ssz.DivideInt2(len(buf), 56, 13)
		if err != nil {
			return err
		}
		s.Committee = make([]*Operator, num)
		for ii := 0; ii < num; ii++ {
			if s.Committee[ii] == nil {
				s.Committee[ii] = new(Operator)
			}
			if err = s.Committee[ii].UnmarshalSSZ(buf[ii*56 : (ii+1)*56]); err != nil {
				return err
			}
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the Share object
func (s *Share) SizeSSZ() (size int) {
	size = 164

	// Field (3) 'Committee'
	size += len(s.Committee) * 56

	return
}

// HashTreeRoot ssz hashes the Share object
func (s *Share) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(s)
}

// HashTreeRootWith ssz hashes the Share object with a hasher
func (s *Share) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'OperatorID'
	hh.PutUint64(uint64(s.OperatorID))

	// Field (1) 'ValidatorPubKey'
	if size := len(s.ValidatorPubKey); size != 48 {
		err = ssz.ErrBytesLengthFn("Share.ValidatorPubKey", size, 48)
		return
	}
	hh.PutBytes(s.ValidatorPubKey)

	// Field (2) 'SharePubKey'
	if size := len(s.SharePubKey); size != 48 {
		err = ssz.ErrBytesLengthFn("Share.SharePubKey", size, 48)
		return
	}
	hh.PutBytes(s.SharePubKey)

	// Field (3) 'Committee'
	{
		subIndx := hh.Index()
		num := uint64(len(s.Committee))
		if num > 13 {
			err = ssz.ErrIncorrectListSize
			return
		}
		for _, elem := range s.Committee {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 13)
	}

	// Field (4) 'DomainType'
	hh.PutBytes(s.DomainType[:])

	// Field (5) 'FeeRecipientAddress'
	hh.PutBytes(s.FeeRecipientAddress[:])

	// Field (6) 'Graffiti'
	if size := len(s.Graffiti); size != 32 {
		err = ssz.ErrBytesLengthFn("Share.Graffiti", size, 32)
		return
	}
	hh.PutBytes(s.Graffiti)

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the Share object
func (s *Share) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(s)
}