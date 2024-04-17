package testingutils

import (
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
)

func UnknownDutyValueCheck() qbft.ProposedValueCheckF {
	return func(data []byte) error {
		return nil
	}
}
