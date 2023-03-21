package spec

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sort"

	consensus "github.com/umbracle/go-eth-consensus"
	"github.com/umbracle/go-eth-consensus/bls"
	"github.com/umbracle/go-eth-consensus/deposit"
)

func ProcessAttestation(state *consensus.BeaconStatePhase0, attestation *consensus.Attestation) error {
	data := attestation.Data

	if data.Target.Epoch != getPreviousEpoch(state) && data.Target.Epoch != getCurrentEpoch(state) {
		return fmt.Errorf("one")
	}
	if data.Target.Epoch != computeEpochAtSlot(data.Slot) {
		return fmt.Errorf("two")
	}

	if !(state.Slot <= data.Slot+Spec.SlotsPerEpoch) {
		return fmt.Errorf("attestation slot is too old")
	}
	if !(data.Slot+Spec.MinAttestationInclusionDelay <= state.Slot) {
		return fmt.Errorf("attestation is too new")
	}

	if data.Index >= getCommitteeCountPerSlot(state, data.Target.Epoch) {
		return fmt.Errorf("ten")
	}

	proposerIndex := getBeaconProposerIndex(state)

	pendingAttestation := &consensus.PendingAttestation{
		Data:            data,
		AggregationBits: attestation.AggregationBits,
		InclusionDelay:  state.Slot - data.Slot,
		ProposerIndex:   proposerIndex,
	}

	if data.Target.Epoch == getCurrentEpoch(state) {
		if *data.Source != *state.CurrentJustifiedCheckpoint {
			return fmt.Errorf("three")
		}
		state.CurrentEpochAttestations = append(state.CurrentEpochAttestations, pendingAttestation)
	} else {
		if *data.Source != *state.PreviousJustifiedCheckpoint {
			return fmt.Errorf("four")
		}
		state.PreviousEpochAttestations = append(state.PreviousEpochAttestations, pendingAttestation)
	}

	indexedAtt, err := getIndexedAttestation(state, attestation)
	if err != nil {
		return err
	}
	if err := isValidIndexedAttestation(state, indexedAtt); err != nil {
		return err
	}
	return nil
}

func getIndexedAttestation(state *consensus.BeaconStatePhase0, attestation *consensus.Attestation) (*consensus.IndexedAttestation, error) {
	attestingIndices, err := getAttestingIndices(state, attestation.Data, attestation.AggregationBits)
	if err != nil {
		return nil, err
	}

	sort.Slice(attestingIndices, func(i, j int) bool {
		return attestingIndices[i] < attestingIndices[j]
	})

	return &consensus.IndexedAttestation{
		AttestationIndices: attestingIndices,
		Data:               attestation.Data,
		Signature:          attestation.Signature,
	}, nil
}

func isValidIndexedAttestation(state *consensus.BeaconStatePhase0, indexedAttestation *consensus.IndexedAttestation) error {
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

	// verify indices are unique
	for i := 1; i < len(indices); i++ {
		if indices[i-1] == indices[i] {
			return fmt.Errorf("duplicated")
		}
	}

	// check if the indices are inside the bounds of the validator set
	if indices[len(indices)-1] >= uint64(len(state.Validators)) {
		return fmt.Errorf("validators out of bounds")
	}

	pubKeys := []*bls.PublicKey{}
	for _, i := range indices {
		pub := new(bls.PublicKey)
		if err := pub.Deserialize(state.Validators[i].Pubkey[:]); err != nil {
			return err
		}
		pubKeys = append(pubKeys, pub)
	}

	domain, err := getDomain(consensus.DomainBeaconAttesterType, state, &indexedAttestation.Data.Target.Epoch)
	if err != nil {
		return err
	}

	root, err := consensus.ComputeSigningRoot(domain, indexedAttestation.Data)
	if err != nil {
		return err
	}

	sig := new(bls.Signature)
	if err := sig.Deserialize(indexedAttestation.Signature[:]); err != nil {
		return err
	}

	ok, err := sig.FastAggregateVerify(pubKeys, root[:])
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("failed to verify")
	}

	return nil
}

func isSlashableAttestationData(d1, d2 *consensus.AttestationData) (bool, error) {
	hash1, err := d1.HashTreeRoot()
	if err != nil {
		return false, fmt.Errorf("unable to hash attestation data 1: %v", err)
	}

	hash2, err := d2.HashTreeRoot()
	if err != nil {
		return false, fmt.Errorf("unable to hash attestation data 2: %v", err)
	}

	return (hash1 != hash2 && d1.Target.Epoch == d2.Target.Epoch) || (d1.Source.Epoch < d2.Source.Epoch && d2.Target.Epoch < d1.Target.Epoch), nil
}

func ProcessAttesterSlashing(state *consensus.BeaconStatePhase0, attesterSlashing *consensus.AttesterSlashing) error {
	att1 := attesterSlashing.Attestation1
	att2 := attesterSlashing.Attestation2

	slashable, err := isSlashableAttestationData(att1.Data, att2.Data)
	if err != nil {
		return fmt.Errorf("unable to determine if attestation data was slashable: %v", err)
	}
	if !slashable {
		return fmt.Errorf("attestation data not slashable: %+v; %+v", att1.Data, att2.Data)
	}

	// Check if attestation 1 is valid
	if err := isValidIndexedAttestation(state, att1); err != nil {
		return fmt.Errorf("failed to validate attestation 1: %v", err)
	}
	// Check if attestation 2 is valid
	if err := isValidIndexedAttestation(state, att2); err != nil {
		return fmt.Errorf("failed to validate attestation 2: %v", err)
	}

	slashedAny := false
	indices := intersection(att1.AttestationIndices, att2.AttestationIndices)

	sort.Slice(indices, func(i, j int) bool {
		return indices[i] < indices[j]
	})

	for _, index := range indices {
		if isSlashableValidator(state.Validators[index], getCurrentEpoch(state)) {
			if err := slashValidator(state, index, nil); err != nil {
				return err
			}
			slashedAny = true
		}
	}

	if !slashedAny {
		return fmt.Errorf("someone should have been slashed?")
	}
	return nil
}

func intersection(a, b []uint64) (res []uint64) {
	found := map[uint64]struct{}{}
	for _, i := range a {
		for _, j := range b {
			if i == j {
				found[i] = struct{}{}
			}
		}
	}

	for i := range found {
		res = append(res, i)
	}
	return
}

func computeProposerIndex(state *consensus.BeaconStatePhase0, indices []uint64, seed [32]byte) uint64 {
	if len(indices) == 0 {
		panic(fmt.Errorf("must have >0 indices"))
	}
	maxRandomByte := uint64(1<<8 - 1)
	i := uint64(0)
	total := uint64(len(indices))
	hash := sha256.New()
	buf := make([]byte, 8)
	for {
		shuffled := computeShuffleIndex(i%total, total, seed)

		candidateIndex := indices[shuffled]
		if candidateIndex >= uint64(len(state.Validators)) {
			panic(fmt.Errorf("candidate index out of range: %d for validator set of length: %d", candidateIndex, len(state.Validators)))
		}
		binary.LittleEndian.PutUint64(buf, i/32)
		input := append(seed[:], buf...)
		hash.Reset()
		hash.Write(input)
		randomByte := uint64(hash.Sum(nil)[i%32])
		effectiveBalance := state.Validators[candidateIndex].EffectiveBalance
		if effectiveBalance*maxRandomByte >= Spec.MaxEffectiveBalance*randomByte {
			return candidateIndex
		}
		i += 1
	}
}

func getEpochAtSlot(slot uint64) uint64 {
	return slot / Spec.SlotsPerEpoch
}

func getBeaconProposerIndex(state *consensus.BeaconStatePhase0) uint64 {
	epoch := getEpochAtSlot(state.Slot)

	hash := sha256.New()
	// Input for the seed hash.
	input := getSeed(state, epoch, consensus.DomainBeaconProposerType)
	slotByteArray := make([]byte, 8)
	binary.LittleEndian.PutUint64(slotByteArray, state.Slot)

	// Add slot to the end of the input.
	inputWithSlot := append(input[:], slotByteArray...)

	// Calculate the hash.
	hash.Write(inputWithSlot)
	seed := hash.Sum(nil)

	indices := getActiveValidatorIndices(state, epoch)

	// Write the seed to an array.
	seedArray := [32]byte{}
	copy(seedArray[:], seed)

	return computeProposerIndex(state, indices, seedArray)
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
	proposerIndex := getBeaconProposerIndex(state)
	if block.ProposerIndex != proposerIndex {
		return fmt.Errorf("incorrect proposer index '%d', expected '%d'", block.ProposerIndex, proposerIndex)
	}

	// Verify that the parent matches
	parentRoot, err := state.LatestBlockHeader.HashTreeRoot()
	if err != nil {
		return err
	}
	if !bytes.Equal(block.ParentRoot[:], parentRoot[:]) {
		return fmt.Errorf("incorrect parent root hash '%s', expected '%s'", hex.EncodeToString(block.ParentRoot[:]), hex.EncodeToString(parentRoot[:]))
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
			// failures in the deposit are tolerated
			return nil
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

	return bytes.Equal(value[:], root[:])
}

func isSlashableValidator(validator *consensus.Validator, epoch uint64) bool {
	return !validator.Slashed && validator.ActivationEpoch <= epoch && epoch < validator.WithdrawableEpoch
}

func decreaseBalance(state *consensus.BeaconStatePhase0, index uint64, delta uint64) {
	if delta > state.Balances[index] {
		state.Balances[index] = 0
	} else {
		state.Balances[index] -= delta
	}
}

func increaseBalance(state *consensus.BeaconStatePhase0, index uint64, delta uint64) {
	state.Balances[index] += delta
}

func slashValidator(state *consensus.BeaconStatePhase0, slashedIndex uint64, whistleblowerIndexPtr *uint64) error {
	epoch := getCurrentEpoch(state)
	if err := initiateValidatorExit(state, slashedIndex); err != nil {
		return err
	}

	validator := state.Validators[slashedIndex]
	validator.Slashed = true
	validator.WithdrawableEpoch = max(validator.WithdrawableEpoch, epoch+Spec.EpochsPerSlashingsVector)

	state.Slashings[epoch%Spec.EpochsPerSlashingsVector] += validator.EffectiveBalance
	decreaseBalance(state, slashedIndex, validator.EffectiveBalance/Spec.MinSlashingPenaltyQuotient)

	// Apply proposer and whistleblower rewards
	proposerIndex := getBeaconProposerIndex(state)

	var whistleblowerIndex uint64
	if whistleblowerIndexPtr != nil {
		whistleblowerIndex = *whistleblowerIndexPtr
	} else {
		whistleblowerIndex = proposerIndex
	}

	whistleblowerReward := validator.EffectiveBalance / Spec.WhistleblowerRewardQuotient
	proposerReward := whistleblowerReward / Spec.ProposerRewardQuotient

	increaseBalance(state, proposerIndex, proposerReward)
	increaseBalance(state, whistleblowerIndex, whistleblowerReward-proposerReward)

	return nil
}

func blsVerify(pubKey []byte, signature []byte, root [32]byte) (bool, error) {
	sig := new(bls.Signature)
	if err := sig.Deserialize(signature); err != nil {
		return false, fmt.Errorf("five: %v", err)
	}
	pub := new(bls.PublicKey)
	if err := pub.Deserialize(pubKey); err != nil {
		return false, fmt.Errorf("six: %v", err)
	}

	ok, err := sig.VerifyByte(pub, root[:])
	if err != nil {
		return false, fmt.Errorf("seven: %v", err)
	}
	return ok, nil
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
	{
		h1, err := header1.HashTreeRoot()
		if err != nil {
			return err
		}
		h2, err := header2.HashTreeRoot()
		if err != nil {
			return err
		}
		if bytes.Equal(h1[:], h2[:]) {
			return fmt.Errorf("three")
		}
	}

	// Verify the proposer is slashable
	if header1.ProposerIndex >= uint64(len(state.Validators)) {
		return fmt.Errorf("four1")
	}
	proposer := state.Validators[header1.ProposerIndex]
	if !isSlashableValidator(proposer, getCurrentEpoch(state)) {
		return fmt.Errorf("four")
	}

	// Verify signatures
	verifySignature := func(signedHeader *consensus.SignedBeaconBlockHeader) error {
		epoch := computeEpochAtSlot(signedHeader.Header.Slot)

		domain, err := getDomain(consensus.DomainBeaconProposerType, state, &epoch)
		if err != nil {
			return err
		}

		root, err := consensus.ComputeSigningRoot(domain, signedHeader.Header)
		if err != nil {
			return err
		}

		ok, err := blsVerify(proposer.Pubkey[:], signedHeader.Signature[:], root)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("failed to verify")
		}

		return nil
	}

	if err := verifySignature(proposerSlashing.Header1); err != nil {
		return err
	}
	if err := verifySignature(proposerSlashing.Header2); err != nil {
		return err
	}

	if err := slashValidator(state, header1.ProposerIndex, nil); err != nil {
		return err
	}
	return nil
}

func computeActivationExitEpoch(epoch uint64) uint64 {
	return epoch + 1 + Spec.MaxSeedLookAhead
}

func getValidatorChurnLimit(state *consensus.BeaconStatePhase0) uint64 {
	activeValidatorIndices := getActiveValidatorIndices(state, getCurrentEpoch(state))

	churnLimit := uint64(len(activeValidatorIndices)) / Spec.ChurnLimitQuotient
	if churnLimit < Spec.MinPerEpochChurnLimit {
		churnLimit = Spec.MinPerEpochChurnLimit
	}
	return churnLimit
}

func initiateValidatorExit(state *consensus.BeaconStatePhase0, index uint64) error {

	// Return if validator already initiated exit
	validator := state.Validators[index]
	if validator.ExitEpoch != farFutureEpoch {
		return nil
	}

	// Compute exit queue epoch
	exitEpochs := []uint64{}
	for _, v := range state.Validators {
		if v.ExitEpoch != farFutureEpoch {
			exitEpochs = append(exitEpochs, v.ExitEpoch)
		}
	}
	exitEpochs = append(exitEpochs, computeActivationExitEpoch(getCurrentEpoch(state)))

	exitQueueEpoch := uint64(0)
	for _, epoch := range exitEpochs {
		if exitQueueEpoch < epoch {
			exitQueueEpoch = epoch
		}
	}

	exitEpochChurn := uint64(0)
	for _, v := range state.Validators {
		if v.ExitEpoch == exitQueueEpoch {
			exitEpochChurn++
		}
	}

	churnLimit := getValidatorChurnLimit(state)
	if exitEpochChurn >= churnLimit {
		exitQueueEpoch++
	}

	// Set validator exit epoch and withdrawable epoch
	validator.ExitEpoch = exitQueueEpoch

	withdrawalEpoch := validator.ExitEpoch + Spec.MinValidatorWithdrawabilityDelay
	if withdrawalEpoch < exitQueueEpoch {
		return fmt.Errorf("overflow epoch")
	}
	validator.WithdrawableEpoch = withdrawalEpoch
	return nil
}

func ProcessVoluntaryExit(state *consensus.BeaconStatePhase0, signedVoluntaryExit *consensus.SignedVoluntaryExit) error {
	voluntaryExit := signedVoluntaryExit.Exit
	if voluntaryExit.ValidatorIndex >= uint64(len(state.Validators)) {
		return fmt.Errorf("bad length")
	}
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
	domain, err := getDomain(consensus.DomainVoluntaryExitType, state, nil)
	if err != nil {
		return fmt.Errorf("three %v", err)
	}
	signingRoot, err := consensus.ComputeSigningRoot(domain, voluntaryExit)
	if err != nil {
		return fmt.Errorf("four %v", err)
	}

	ok, err := blsVerify(validator.Pubkey[:], signedVoluntaryExit.Signature[:], signingRoot)
	if err != nil {
		return fmt.Errorf("seven: %v", err)
	}
	if !ok {
		return fmt.Errorf("failed to validate")
	}

	// Initiate exit
	if err := initiateValidatorExit(state, voluntaryExit.ValidatorIndex); err != nil {
		return err
	}

	return nil
}

func getDomain(domain consensus.Domain, state *consensus.BeaconStatePhase0, epoch *uint64) ([32]byte, error) {
	var forkVersion [4]byte

	curEpoch := getCurrentEpoch(state)
	if epoch != nil {
		curEpoch = *epoch
	}

	if curEpoch < state.Fork.Epoch {
		forkVersion = state.Fork.PreviousVersion
	} else {
		forkVersion = state.Fork.CurrentVersion
	}
	return consensus.ComputeDomain(domain, forkVersion, state.GenesisValidatorsRoot)
}
