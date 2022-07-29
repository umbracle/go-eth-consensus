package deposit

import (
	"encoding/hex"

	ssz "github.com/ferranbt/fastssz"
	"github.com/umbracle/ethgo/abi"
	consensus "github.com/umbracle/go-eth-consensus"
	"github.com/umbracle/go-eth-consensus/bls"
)

const MinGweiAmount = uint64(320)

// DepositEvent is the eth2 deposit event
var DepositEvent = abi.MustNewEvent(`event DepositEvent(
	bytes pubkey,
	bytes whitdrawalcred,
	bytes amount,
	bytes signature,
	bytes index
)`)

var depositDomain [32]byte

func init() {
	// the domain for the deposit signing is hardcoded
	buf, _ := hex.DecodeString("03000000f5a5fd42d16a20302798ef6ed309979b43003d2320d9f0e8ea9831a9")
	copy(depositDomain[:], buf)
}

func Input(depositKey *bls.Key, withdrawalKey *bls.Key, amountInGwei uint64) (*consensus.DepositData, error) {
	// withdrawalCredentialsHash forms a 32 byte hash of the withdrawal public address.
	//   withdrawal_credentials[:1] == BLS_WITHDRAWAL_PREFIX_BYTE
	//   withdrawal_credentials[1:] == hash(withdrawal_pubkey)[1:]
	// TODO

	unsignedMsgRoot, err := ssz.HashWithDefaultHasher(&consensus.DepositMessage{
		Pubkey:                depositKey.Pub.Serialize(),
		Amount:                amountInGwei,
		WithdrawalCredentials: [32]byte{},
	})
	if err != nil {
		return nil, err
	}

	rootToSign, err := ssz.HashWithDefaultHasher(&consensus.SigningData{
		ObjectRoot: unsignedMsgRoot,
		Domain:     depositDomain,
	})
	if err != nil {
		return nil, err
	}

	signature, err := depositKey.Sign(rootToSign)
	if err != nil {
		return nil, err
	}

	msg := &consensus.DepositData{
		Pubkey:                depositKey.Pub.Serialize(),
		Amount:                amountInGwei,
		WithdrawalCredentials: [32]byte{},
		Signature:             signature,
	}
	root, err := msg.HashTreeRoot()
	if err != nil {
		return nil, err
	}
	msg.Root = root
	return msg, nil
}

func signingData(obj ssz.HashRoot) ([32]byte, error) {
	unsignedMsgRoot, err := ssz.HashWithDefaultHasher(obj)
	if err != nil {
		return [32]byte{}, err
	}

	root, err := ssz.HashWithDefaultHasher(&consensus.SigningData{
		ObjectRoot: unsignedMsgRoot,
		Domain:     depositDomain,
	})
	if err != nil {
		return [32]byte{}, err
	}
	return root, nil
}
