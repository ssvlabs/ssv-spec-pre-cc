package messages

import (
	"encoding/hex"

	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft/spectest/tests"
	"github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

// MarshalJustificationsWithFullData tests marshalling justifications with full data (should omit it)
func MarshalJustificationsWithFullData() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*qbft.SignedMessage{
		testingutils.TestingRoundChangeMessageWithRound(ks.Shares[1], 1, 2),
	}

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.TestingPrepareMessage(ks.Shares[1], types.OperatorID(1)),
	}

	msg := testingutils.TestingProposalMessageWithParams(
		ks.Shares[1], types.OperatorID(1), 2, qbft.FirstHeight, testingutils.TestingQBFTRootData,
		testingutils.MarshalJustifications(rcMsgs), testingutils.MarshalJustifications(prepareMsgs))

	r, err := hex.DecodeString("fddc24432bfcc56695474576b1c70aed011f18bfa1ad2c10c85772f8c588e412")
	if err != nil {
		panic(err)
	}

	b, err := hex.DecodeString("901603cd1e6e9c2a4fd690039a3f0731e17a64d1a68ab883dbff61688822054870e300fa9f21b30ff315f0851376758100ff0919db5d6dc674ce71577426c724d8935adba2c624ad2323fac3510c2dd117eb49f6a3a6b4171fb749a67aa0772e6c000000740000005402000001000000000000000000000000000000000000000000000002000000000000004c000000be956fb7df4ef37531682d588320084fc914c3f0fed335263e5b44062e6c29b4000000000000000050000000180100000102030404000000b77c035c1d1d9c6c7cc22810c50a85c9e560ad791509055061dfe07d403edaa1304161466d49f64844b6e5cf4b09709f069debfd91438f97d414a4f64cdcb4f8cf1e9703f8da141be9f74509087e2f84492314b0341966c4b8fc16d931f6ba086c00000074000000c400000001000000000000000300000000000000000000000000000002000000000000004c000000be956fb7df4ef37531682d588320084fc914c3f0fed335263e5b44062e6c29b40000000000000000500000005000000001020304040000008129e6862a5120bd085e1936b4efb5a55fc7d19c0d0fda0e9ec576d18abd4a17ab3a033f5296b74c5fdaf85cb7b3da3201b63feca76b883613e3b1ca137e763a342e3b1dddbce016f8ca3cbce32c8b125dd8c25a7639819c20b539e9e7c6c5796c00000074000000c400000001000000000000000100000000000000000000000000000001000000000000004c000000be956fb7df4ef37531682d588320084fc914c3f0fed335263e5b44062e6c29b40000000000000000500000005000000001020304010203040506070809010203040506070809010203040506070809")
	if err != nil {
		panic(err)
	}

	return &tests.MsgSpecTest{
		Name: "marshal justifications with full data",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		EncodedMessages: [][]byte{
			b,
		},
		ExpectedRoots: [][32]byte{
			*(*[32]byte)(r),
		},
	}
}
