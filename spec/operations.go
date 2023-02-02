package spec

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"sort"

	consensus "github.com/umbracle/go-eth-consensus"
	"github.com/umbracle/go-eth-consensus/deposit"
)

func isValidIndexAttestation(pre *consensus.BeaconStatePhase0, indexedAttestation *consensus.IndexedAttestation) error {
	indices := indexedAttestation.AttestationIndices

	// the attestation cannot be empty
	if len(indices) == 0 {
		return fmt.Errorf("empty")
	}

	// the indices have to be sorted
	isSorted := sort.SliceIsSorted(indices, func(i, j int) bool {
		return indices[i] < indices[j]
	})
	if !isSorted {
		return fmt.Errorf("indices not sorted")
	}

	return nil
}

func ProcessAttesterSlashing(pre *consensus.BeaconStatePhase0, attesterSlashing *consensus.AttesterSlashing) error {
	att1 := attesterSlashing.Attestation1
	att2 := attesterSlashing.Attestation2

	// Check if ``att1.data`` and ``att2.data`` are slashable according to Casper FFG rules:
	// 1. Double vote.
	if att1.Data != att2.Data && att1.Data.Target.Epoch == att2.Data.Target.Epoch {
		return fmt.Errorf("double vote")
	}
	// 2. Surround vote.
	if att1.Data.Source.Epoch < att2.Data.Source.Epoch && att2.Data.Target.Epoch < att1.Data.Target.Epoch {
		return fmt.Errorf("surround vote")
	}

	// Check if attestation 1 is valid
	if err := isValidIndexAttestation(pre, att1); err != nil {
		return fmt.Errorf("failed to validate attestation 1: %v", err)
	}
	// Check if attestation 2 is valid
	if err := isValidIndexAttestation(pre, att2); err != nil {
		return fmt.Errorf("failed to validate attestation 2: %v", err)
	}

	return nil
}

func getBeaconProposerIndex(state *consensus.BeaconStatePhase0) (uint64, error) {
	epoch := getCurrentEpoch(state)
	fmt.Println(epoch)

	return 0, nil
}

func ProcessBlockHeader(state *consensus.BeaconStatePhase0, block *consensus.BeaconBlockPhase0) error {
	// Verify that the slots match
	if block.Slot != state.Slot {
		return fmt.Errorf("slot mismatch: %d, %d", block.Slot, state.Slot)
	}

	// Verify that the block is newer than latest block header
	if block.Slot <= state.LatestBlockHeader.Slot {
		return fmt.Errorf("block slot %d is older than latest block header %d", state.LatestBlockHeader.Slot, block.Slot)
	}

	// Verify that proposer index is the correct index
	proposerIndex, err := getBeaconProposerIndex(state)
	if err != nil {
		return err
	}
	if block.ProposerIndex != proposerIndex {
		return fmt.Errorf("incorrect proposer index '%d', expected '%d'", block.ProposerIndex, proposerIndex)
	}

	// Verify that the parent matches
	parentRoot, err := state.LatestBlockHeader.HashTreeRoot()
	if err != nil {
		return err
	}
	if !bytes.Equal(block.ParentRoot[:], parentRoot[:]) {
		return fmt.Errorf("incorrect parent root hash '%s', expected '%s'", block.ParentRoot[:], parentRoot[:])
	}

	// Cache current block as the new latest block
	bodyRoot, err := block.Body.HashTreeRoot()
	if err != nil {
		return err
	}
	state.LatestBlockHeader = &consensus.BeaconBlockHeader{
		Slot:          block.Slot,
		ProposerIndex: block.ProposerIndex,
		ParentRoot:    block.ParentRoot,
		BodyRoot:      bodyRoot,
	}

	// Verify proposer is not slashed
	if state.Validators[block.ProposerIndex].Slashed {
		return fmt.Errorf("proposer is slashed")
	}

	return nil
}

const (
	depositContractTreeDepth = 32

	farFutureEpoch = 18446744073709551615 // 2**64-1
)

func ProcessDeposit(state *consensus.BeaconStatePhase0, depositObj *consensus.Deposit) error {
	// Verify the Merkle branch
	depositRoot, err := depositObj.Data.HashTreeRoot()
	if err != nil {
		return err
	}
	if !isValidMerkleBranch(depositRoot, depositObj.Proof, depositContractTreeDepth+1, state.Eth1DepositIndex, state.Eth1Data.DepositRoot) {
		return fmt.Errorf("bad merkle root")
	}

	// Deposits must be processed in order
	state.Eth1DepositIndex++

	pubKey := depositObj.Data.Pubkey
	amount := depositObj.Data.Amount

	indx, ok := isInValidatorSet(state, pubKey)
	if !ok {
		// Verify the deposit signature (proof of possession) which is not checked by the deposit contract
		if err := deposit.Verify(depositObj.Data); err != nil {
			return err
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

	return nil
}

func isInValidatorSet(state *consensus.BeaconStatePhase0, pubKey [48]byte) (uint64, bool) {
	for indx, val := range state.Validators {
		if bytes.Equal(val.Pubkey[:], pubKey[:]) {
			return uint64(indx), true
		}
	}
	return 0, false
}

func isValidMerkleBranch(leaf [32]byte, proof [33][32]byte, depth, index uint64, root [32]byte) bool {
	value := leaf

	for i := uint64(0); i < depth; i++ {
		if (index>>i)&1 == 1 {
			value = sha256.Sum256(append(proof[i][:], value[:]...))
		} else {
			value = sha256.Sum256(append(value[:], proof[i][:]...))
		}
	}

	return !bytes.Equal(value[:], root[:])
}

func ProcessProposerSlashing(state *consensus.BeaconStatePhase0, proposerSlashing *consensus.ProposerSlashing) error {
	header1 := proposerSlashing.Header1.Header
	header2 := proposerSlashing.Header2.Header

	// Verify header slots match
	if header1.Slot != header2.Slot {
		return fmt.Errorf("one")
	}

	// Verify header proposer indices match
	if header1.ProposerIndex != header2.ProposerIndex {
		return fmt.Errorf("two")
	}

	// Verify the headers are different
	if header1 != header2 {
		return fmt.Errorf("three")
	}

	// Verify the proposer is slashable

	// Verify signatures

	return nil
}

func ProcessVoluntaryExit(state *consensus.BeaconStatePhase0, signedVoluntaryExit *consensus.SignedVoluntaryExit) error {
	voluntaryExit := signedVoluntaryExit.Exit
	validator := state.Validators[voluntaryExit.ValidatorIndex]

	// Verify the validator is active
	if !isActiveValidator(validator, getCurrentEpoch(state)) {
		return fmt.Errorf("the validator is not active")
	}

	// Verify exit has not been initiated
	if validator.ExitEpoch != farFutureEpoch {
		return fmt.Errorf("exit has already been initialized")
	}

	// Exits must specify an epoch when they become valid; they are not valid before then
	if getCurrentEpoch(state) < voluntaryExit.Epoch {
		return fmt.Errorf("one")
	}

	// Verify the validator has been active long enough
	if getCurrentEpoch(state) < validator.ActivationEpoch+Spec.ShardCommiteePeriod {
		return fmt.Errorf("two")
	}

	// Verify signature

	// Initiate exit

	return nil
}
