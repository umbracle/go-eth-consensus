package spec

import (
	"bytes"
	"sort"

	consensus "github.com/umbracle/go-eth-consensus"
)

func processEffectiveBalanceUpdates(state *consensus.BeaconStatePhase0) error {
	for indx, validator := range state.Validators {
		balance := state.Balances[indx]

		hysteresisIncrement := Spec.EffectiveBalanceIncrement / Spec.HysteresisQuotient
		downwardThreshold := hysteresisIncrement * Spec.HysteresisDownwardMultiplier
		upwardThreshold := hysteresisIncrement * Spec.HysteresisUpwardMultiplier

		if balance+downwardThreshold < validator.EffectiveBalance || validator.EffectiveBalance+upwardThreshold < balance {
			validator.EffectiveBalance = min(balance-balance%Spec.EffectiveBalanceIncrement, Spec.MaxEffectiveBalance)
		}
	}
	return nil
}

func processEth1DataReset(state *consensus.BeaconStatePhase0) error {
	nextEpoch := getCurrentEpoch(state) + 1

	if nextEpoch%Spec.EpochsPerEth1VotingPeriod == 0 {
		state.Eth1DataVotes = []*consensus.Eth1Data{}
	}
	return nil
}

func processHistoricalRootsUpdate(state *consensus.BeaconStatePhase0) error {
	nextEpoch := getCurrentEpoch(state) + 1

	if nextEpoch%(Spec.SlotsPerHistoricalRoot/Spec.SlotsPerEpoch) == 0 {
		historicalBatch := consensus.HistoricalBatch{
			BlockRoots: state.BlockRoots,
			StateRoots: state.StateRoots,
		}
		root, err := historicalBatch.HashTreeRoot()
		if err != nil {
			return err
		}
		state.HistoricalRoots = append(state.HistoricalRoots, root)
	}
	return nil
}

func computeStartSlotAtEpoch(epoch uint64) uint64 {
	return epoch * Spec.SlotsPerEpoch
}

func getBlockRootAtSlot(state *consensus.BeaconStatePhase0, slot uint64) [32]byte {
	// Return the block root at a recent ``slot``.
	return state.BlockRoots[slot%Spec.SlotsPerHistoricalRoot]
}

func getBlockRoot(state *consensus.BeaconStatePhase0, epoch uint64) [32]byte {
	return getBlockRootAtSlot(state, computeStartSlotAtEpoch(epoch))
}

func getMatchingTargetAttestations(state *consensus.BeaconStatePhase0, epoch uint64) []*consensus.PendingAttestation {
	root := getBlockRoot(state, epoch)

	res := []*consensus.PendingAttestation{}
	for _, a := range getMatchingSourceAttestations(state, epoch) {
		if bytes.Equal(a.Data.Target.Root[:], root[:]) {
			res = append(res, a)
		}
	}

	return res
}

func getAttestingBalance(state *consensus.BeaconStatePhase0, attestations []*consensus.PendingAttestation) uint64 {
	// Return the combined effective balance of the set of unslashed validators participating in ``attestations``.
	// Note: ``get_total_balance`` returns ``EFFECTIVE_BALANCE_INCREMENT`` Gwei minimum to avoid divisions by zero
	indices, err := getUnslashedAttestingIndices(state, attestations)
	if err != nil {
		panic(err)
	}
	return getTotalBalance(state, indices)
}

func processJustificationAndFinalization(state *consensus.BeaconStatePhase0) error {
	// Initial FFG checkpoint values have a `0x00` stub for `root`.
	// Skip FFG updates in the first two epochs to avoid corner cases that might result in modifying this stub.
	if getCurrentEpoch(state) <= Spec.GenesisEpoch+1 {
		return nil
	}

	previousAttestations := getMatchingSourceAttestations(state, getPreviousEpoch(state))
	currentAttestations := getMatchingTargetAttestations(state, getCurrentEpoch(state))
	totalActiveBalance := getTotalActiveBalance(state)

	previousTargetBalalnce := getAttestingBalance(state, previousAttestations)
	currentTargetBalance := getAttestingBalance(state, currentAttestations)

	weighJustificationAndFinalization(state, totalActiveBalance, previousTargetBalalnce, currentTargetBalance)
	return nil
}

func weighJustificationAndFinalization(state *consensus.BeaconStatePhase0, totalActiveBalance uint64, previousEpochTargetBalance uint64, currentEpochTargetBalance uint64) {
	previousEpoch := getPreviousEpoch(state)
	currentEpoch := getCurrentEpoch(state)

	oldPreviousJustifiedCheckpoint := state.PreviousJustifiedCheckpoint
	oldCurrentJustifiedCheckpoint := state.CurrentJustifiedCheckpoint

	// Process justifications
	state.PreviousJustifiedCheckpoint = state.CurrentJustifiedCheckpoint
	state.JustificationBits[0] = (state.JustificationBits[0] << 1) & 0x0f

	if previousEpochTargetBalance*3 >= totalActiveBalance*2 {
		state.CurrentJustifiedCheckpoint = &consensus.Checkpoint{
			Epoch: previousEpoch,
			Root:  getBlockRoot(state, previousEpoch),
		}

		state.JustificationBits[0] |= 1 << 1
	}

	if currentEpochTargetBalance*3 >= totalActiveBalance*2 {
		state.CurrentJustifiedCheckpoint = &consensus.Checkpoint{
			Epoch: currentEpoch,
			Root:  getBlockRoot(state, currentEpoch),
		}

		state.JustificationBits[0] |= 1 << 0
	}

	// Process finalizations
	bits := state.JustificationBits[0]

	// the 2nd/3rd/4th most recent epochs are justified, the 2nd using the 4th as source
	if bits&0x0E == 0x0E && oldPreviousJustifiedCheckpoint.Epoch+3 == currentEpoch {
		state.FinalizedCheckpoint = oldPreviousJustifiedCheckpoint
	}
	// the 2nd/3rd most recent epochs are justified, the 2nd using the 3rd as source
	if bits&0x06 == 0x06 && oldPreviousJustifiedCheckpoint.Epoch+2 == currentEpoch {
		state.FinalizedCheckpoint = oldPreviousJustifiedCheckpoint
	}
	// the 1st/2nd/3rd most recent epochs are justified, the 1st using the 3rd as source
	if bits&0x07 == 0x07 && oldCurrentJustifiedCheckpoint.Epoch+2 == currentEpoch {
		state.FinalizedCheckpoint = oldCurrentJustifiedCheckpoint
	}
	// the 1st/2nd most recent epochs are justified, the 1st using the 2nd as source
	if bits&0x03 == 0x03 && oldCurrentJustifiedCheckpoint.Epoch+1 == currentEpoch {
		state.FinalizedCheckpoint = oldCurrentJustifiedCheckpoint
	}
}

func processParticipationRecordUpdates(state *consensus.BeaconStatePhase0) error {
	state.PreviousEpochAttestations = state.CurrentEpochAttestations
	state.CurrentEpochAttestations = []*consensus.PendingAttestation{}
	return nil
}

func processRandaoMixesReset(state *consensus.BeaconStatePhase0) error {
	currentEpoch := getCurrentEpoch(state)
	nextEpoch := currentEpoch + 1
	state.RandaoMixes[nextEpoch%Spec.EpochsPerHistoricalVector] = getRandaoMix(state, currentEpoch)
	return nil
}

func isElegibleForActivationQueue(validator *consensus.Validator) bool {
	return validator.ActivationEligibilityEpoch == farFutureEpoch && validator.EffectiveBalance == Spec.MaxEffectiveBalance
}

func isElegibleForActivation(state *consensus.BeaconStatePhase0, validator *consensus.Validator) bool {
	return validator.ActivationEligibilityEpoch <= state.FinalizedCheckpoint.Epoch && validator.ActivationEpoch == farFutureEpoch
}

func processRegistryUpdates(state *consensus.BeaconStatePhase0) error {
	// Process activation eligibility and ejections
	for indx, validator := range state.Validators {
		if isElegibleForActivationQueue(validator) {
			validator.ActivationEligibilityEpoch = getCurrentEpoch(state) + 1
		}

		if isActiveValidator(validator, getCurrentEpoch(state)) && validator.EffectiveBalance <= Spec.EjectionBalance {
			if err := initiateValidatorExit(state, uint64(indx)); err != nil {
				return err
			}
		}
	}

	// Queue validators eligible for activation and not yet dequeued for activation
	activationQueue := []uint64{}
	for indx, validator := range state.Validators {
		if isElegibleForActivation(state, validator) {
			activationQueue = append(activationQueue, uint64(indx))
		}
	}

	// Order by the sequence of activation_eligibility_epoch setting and then index
	sort.Slice(activationQueue, func(i, j int) bool {
		if state.Validators[i].ActivationEligibilityEpoch == state.Validators[j].ActivationEligibilityEpoch {
			// Order by index
			return i < j
		}
		// Order by ActivationEligibilityEpoch
		return state.Validators[i].ActivationEligibilityEpoch < state.Validators[j].ActivationEligibilityEpoch
	})

	churnLimit := min(uint64(len(activationQueue)), getValidatorChurnLimit(state))

	// Dequeued validators for activation up to churn limit
	for _, indx := range activationQueue[:churnLimit] {
		validator := state.Validators[indx]
		validator.ActivationEpoch = computeActivationExitEpoch(getCurrentEpoch(state))
	}
	return nil
}

// getInclusionDelayDeltas returns proposer and inclusion delay micro-rewards/penalties for each validator.
func getInclusionDelayDeltas(state *consensus.BeaconStatePhase0) ([]uint64, []uint64) {
	rewards := make([]uint64, len(state.Validators))

	matchingSourceAttestations := getMatchingSourceAttestations(state, getPreviousEpoch(state))

	unslashedAttIndex, err := getUnslashedAttestingIndices(state, matchingSourceAttestations)
	if err != nil {
		panic(err)
	}

	for _, index := range unslashedAttIndex {
		var attestation *consensus.PendingAttestation
		for _, a := range matchingSourceAttestations {
			attIndex, err := getAttestingIndices(state, a.Data, a.AggregationBits)
			if err != nil {
				panic(err)
			}
			if contains(attIndex, index) {
				if attestation != nil {
					if attestation.InclusionDelay < a.InclusionDelay {
						continue
					}
				}
				attestation = a
			}
		}

		rewards[attestation.ProposerIndex] += getProposerReward(state, index)
		maxAttesterReward := getBaseReward(state, index) - getProposerReward(state, index)
		rewards[index] += maxAttesterReward / attestation.InclusionDelay
	}

	// no penalties associated with inclusion delay
	penalties := make([]uint64, len(state.Validators))

	return rewards, penalties
}

func contains(a []uint64, b uint64) bool {
	for _, i := range a {
		if i == b {
			return true
		}
	}
	return false
}

// getInactivityPenaltyDeltas return inactivity reward/penalty deltas for each validator.
func getInactivityPenaltyDeltas(state *consensus.BeaconStatePhase0) ([]uint64, []uint64) {
	penalties := make([]uint64, len(state.Validators))

	if isInInactivityLeak(state) {
		matchingTargetAttestations := getMatchingTargetAttestations(state, getPreviousEpoch(state))
		matchingTargetAttestingIndices, err := getUnslashedAttestingIndices(state, matchingTargetAttestations)
		if err != nil {
			panic(err)
		}

		for _, index := range getElegibleValidatorIndices(state) {
			// If validator is performing optimally this cancels all rewards for a neutral balance
			baseReward := getBaseReward(state, index)
			penalties[index] += Spec.BaseRewardsPerEpoch*baseReward - getProposerReward(state, index)

			if !contains(matchingTargetAttestingIndices, index) {
				effectiveBalance := state.Validators[index].EffectiveBalance
				penalties[index] += effectiveBalance * getFinalityDelay(state) / Spec.InactivityPenaltyQuotient
			}
		}
	}

	// No rewards associated with inactivity penalties
	rewards := make([]uint64, len(state.Validators))

	return rewards, penalties
}

func getProposerReward(state *consensus.BeaconStatePhase0, attestingIndex uint64) uint64 {
	return getBaseReward(state, attestingIndex) / Spec.ProposerRewardQuotient
}

func getAttestationDeltas(state *consensus.BeaconStatePhase0) ([]uint64, []uint64) {
	// Return attestation reward/penalty deltas for each validator.
	sourceRewards, sourcePenalties := getSourceDeltas(state)
	targetRewards, targetPenalties := getTargetDeltas(state)
	headRewards, headPenalties := getHeadDeltas(state)
	inclusionDelayRewards, _ := getInclusionDelayDeltas(state)
	_, inactivityPenalties := getInactivityPenaltyDeltas(state)

	penalties := make([]uint64, len(state.Validators))
	rewards := make([]uint64, len(state.Validators))

	for i := 0; i < len(state.Validators); i++ {
		rewards[i] = sourceRewards[i] + targetRewards[i] + headRewards[i] + inclusionDelayRewards[i]
	}

	for i := 0; i < len(state.Validators); i++ {
		penalties[i] = sourcePenalties[i] + targetPenalties[i] + headPenalties[i] + inactivityPenalties[i]
	}

	return rewards, penalties
}

func processRewardsAndPenalties(state *consensus.BeaconStatePhase0) error {
	// No rewards are applied at the end of `GENESIS_EPOCH` because rewards are for work done in the previous epoch
	if getCurrentEpoch(state) == Spec.GenesisEpoch {
		return nil
	}

	rewards, penalties := getAttestationDeltas(state)
	for indx := range state.Validators {
		increaseBalance(state, uint64(indx), rewards[indx])
		decreaseBalance(state, uint64(indx), penalties[indx])
	}
	return nil
}

func sum(i []uint64) (res uint64) {
	for _, j := range i {
		res += j
	}
	return
}

func processSlashings(state *consensus.BeaconStatePhase0) error {
	epoch := getCurrentEpoch(state)

	totalBalance := getTotalActiveBalance(state)
	adjustedTotalSlashingBalance := min(sum(state.Slashings)*Spec.ProportionalSlashingsMultiplier, totalBalance)

	for index, validator := range state.Validators {
		if validator.Slashed && epoch+Spec.EpochsPerSlashingsVector/2 == validator.WithdrawableEpoch {
			increment := Spec.EffectiveBalanceIncrement
			penaltyNumerator := (validator.EffectiveBalance / increment) * adjustedTotalSlashingBalance
			penalty := (penaltyNumerator / totalBalance) * increment
			decreaseBalance(state, uint64(index), penalty)
		}
	}
	return nil
}

func processSlashingsReset(state *consensus.BeaconStatePhase0) error {
	nextEpoch := getCurrentEpoch(state) + 1
	state.Slashings[nextEpoch%Spec.EpochsPerSlashingsVector] = 0
	return nil
}
