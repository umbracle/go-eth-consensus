package spec

import (
	"bytes"

	consensus "github.com/umbracle/go-eth-consensus"
	"github.com/umbracle/go-eth-consensus/deposit"
)

const (
	farFutureEpoch = 2*64 - 1
)

func ProcessDeposit(state *consensus.BeaconStatePhase0, depositObj *consensus.Deposit) {
	// Verify the Merkle branch

	// Deposits must be processed in order
	state.Eth1DepositIndex++

	pubKey := depositObj.Data.Pubkey
	amount := depositObj.Data.Amount

	indx, ok := isInValidatorSet(state, pubKey)
	if !ok {
		// Verify the deposit signature (proof of possession) which is not checked by the deposit contract
		if err := deposit.Verify(depositObj.Data); err != nil {
			return
		}

		effectiveBalance := amount - amount%Spec.EffectiveBalanceIncrement
		if effectiveBalance > Spec.MaxEffectiveBalance {
			effectiveBalance = Spec.MaxEffectiveBalance
		}

		val := &consensus.Validator{
			Pubkey:                     pubKey,
			EffectiveBalance:           effectiveBalance,
			WithdrawalCredentials:      depositObj.Data.WithdrawalCredentials,
			ActivationEligibilityEpoch: farFutureEpoch,
			ActivationEpoch:            farFutureEpoch,
			ExitEpoch:                  farFutureEpoch,
			WithdrawableEpoch:          farFutureEpoch,
		}

		// Add validator and balance entries
		state.Validators = append(state.Validators, val)
		state.Balances = append(state.Balances, amount)
	} else {
		// increase balance by deposit amount
		state.Balances[indx] += amount
	}
}

func isInValidatorSet(state *consensus.BeaconStatePhase0, pubKey [48]byte) (uint64, bool) {
	for indx, val := range state.Validators {
		if bytes.Equal(val.Pubkey[:], pubKey[:]) {
			return uint64(indx), true
		}
	}
	return 0, false
}
