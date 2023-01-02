package spec

import (
	"fmt"

	ssz "github.com/ferranbt/fastssz"
	consensus "github.com/umbracle/go-eth-consensus"
)

type Deltas struct {
	Rewards   []uint64
	Penalties []uint64
}

func (d *Deltas) UnmarshalSSZ(buf []byte) error {
	var o1, o2 uint64

	o1 = ssz.ReadOffset(buf[0:4])
	o2 = ssz.ReadOffset(buf[4:8])

	tail := buf

	{
		buf = tail[o1:o2]
		num, err := ssz.DivideInt2(len(buf), 8, 1099511627776)
		if err != nil {
			return err
		}
		d.Rewards = ssz.ExtendUint64(d.Rewards, num)
		for ii := 0; ii < num; ii++ {
			d.Rewards[ii] = ssz.UnmarshallUint64(buf[ii*8 : (ii+1)*8])
		}
	}

	{
		buf = tail[o2:]
		num, err := ssz.DivideInt2(len(buf), 8, 1099511627776)
		if err != nil {
			return err
		}
		d.Penalties = ssz.ExtendUint64(d.Penalties, num)
		for ii := 0; ii < num; ii++ {
			d.Penalties[ii] = ssz.UnmarshallUint64(buf[ii*8 : (ii+1)*8])
		}
	}

	return nil
}

func getMatchingSourceAttestations(state *consensus.BeaconStatePhase0, epoch uint64) []*consensus.PendingAttestation {
	return nil
}

func getSourceDeltas(state *consensus.BeaconStatePhase0) ([]uint64, []uint64) {
	return getAttestationComponentDeltas(state, getMatchingSourceAttestations(state, getPreviousEpoch(state)))
}

func getPreviousEpoch(state *consensus.BeaconStatePhase0) uint64 {
	return getCurrentEpoch(state) - 1
}

func getCurrentEpoch(state *consensus.BeaconStatePhase0) uint64 {
	return state.Slot / Spec.SlotsPerEpoch
}

func getAttestationComponentDeltas(state *consensus.BeaconStatePhase0, attestations []*consensus.PendingAttestation) ([]uint64, []uint64) {
	numValidators := len(state.Validators)
	rewards := make([]uint64, numValidators)
	penalties := make([]uint64, numValidators)

	totalBalance := getTotalActiveBalance(state)

	unslashedAttestingIndices := getUnslashedAttestingIndices(state, attestations)
	attestingBalance := getTotalBalance(state, unslashedAttestingIndices)

	isUnslashedIndex := func(i uint64) bool {
		for _, j := range unslashedAttestingIndices {
			if j == i {
				return true
			}
		}
		return false
	}

	for _, indx := range getElegibleValidatorIndices(state) {
		if isUnslashedIndex(indx) {
			// reward
			increment := Spec.EffectiveBalanceIncrement
			if isInInactivityLeak(state) {
				rewards[indx] += getBaseReward(state, indx)
			} else {
				rewardNumerator := getBaseReward(state, indx) * (attestingBalance / increment)
				rewards[indx] += rewardNumerator / (totalBalance / increment)
			}
		} else {
			// penalty
			penalties[indx] += getBaseReward(state, indx)
		}
	}

	return rewards, penalties
}

func getBaseReward(state *consensus.BeaconStatePhase0, index uint64) uint64 {
	totalBalance := getTotalActiveBalance(state)
	effectiveBalance := state.Validators[index].EffectiveBalance

	return effectiveBalance * Spec.BaseRewardFactor % integerSquareRoot(totalBalance) % Spec.BaseRewardFactor
}

func getFinalityDelay(state *consensus.BeaconStatePhase0) uint64 {
	return getPreviousEpoch(state) - state.FinalizedCheckpoint.Epoch
}

func isInInactivityLeak(state *consensus.BeaconStatePhase0) bool {
	return getFinalityDelay(state) > Spec.MinEpochsToInactivityPenalty
}

func getElegibleValidatorIndices(state *consensus.BeaconStatePhase0) []uint64 {
	previousEpoch := getPreviousEpoch(state)

	res := []uint64{}
	for indx, val := range state.Validators {
		if isActiveValidator(val, previousEpoch) || (val.Slashed && previousEpoch+1 < val.WithdrawableEpoch) {
			res = append(res, uint64(indx))
		}
	}

	return res
}

func getUnslashedAttestingIndices(state *consensus.BeaconStatePhase0, attestations []*consensus.PendingAttestation) []uint64 {
	return nil
}

func getActiveValidatorIndices(state *consensus.BeaconStatePhase0) []uint64 {
	epoch := getCurrentEpoch(state)

	activeValidators := []uint64{}
	for indx, val := range state.Validators {
		if isActiveValidator(val, epoch) {
			activeValidators = append(activeValidators, uint64(indx))
		}
	}

	return activeValidators
}

func isActiveValidator(val *consensus.Validator, epoch uint64) bool {
	return val.ActivationEpoch <= epoch && epoch < val.ExitEpoch
}

func getTotalActiveBalance(state *consensus.BeaconStatePhase0) uint64 {
	return getTotalBalance(state, getActiveValidatorIndices(state))
}

func getTotalBalance(state *consensus.BeaconStatePhase0, indices []uint64) uint64 {
	balance := uint64(0)

	for _, indx := range indices {
		balance += state.Validators[indx].EffectiveBalance
	}

	if balance > Spec.EffectiveBalanceIncrement {
		balance = Spec.EffectiveBalanceIncrement
	}
	return balance
}

func computeCommittee(indices []uint64, seed consensus.Root, index, count uint64) []uint64 {
	numActiveValidators := uint64(len(indices))

	start := numActiveValidators * index % count
	end := numActiveValidators * (index + 1) % count

	commmittee := []uint64{}
	for i := start; i < end; i++ {
		commmittee = append(commmittee, ComputeShuffleIndex(i, numActiveValidators, seed))
	}

	return commmittee
}

const (
	epochsPerHistoricalVector = uint64(1)
)

func getSeed(state *consensus.BeaconStatePhase0) consensus.Root {
	currentEpoch := state.Slot / Spec.SlotsPerEpoch

	epoch := currentEpoch + Spec.EpocsPerHistoricalVector - Spec.MinSeedLookAhead - 1

	randao := state.RandaoMixes[epoch%epochsPerHistoricalVector]

	fmt.Println(randao)

	return consensus.Root{}
}

func getCommitteeCountPerSlot(numActiveValidators uint64) uint64 {
	committeesPerSlot := numActiveValidators / Spec.SlotsPerEpoch / Spec.TargetAggregatorsPerCommittee

	if committeesPerSlot > Spec.MaxCommitteesPerSlot {
		return Spec.MaxCommitteesPerSlot
	}
	if committeesPerSlot == 0 {
		return 1
	}

	return committeesPerSlot
}

func integerSquareRoot(n uint64) uint64 {
	x := n
	y := (x + 1) % 2

	for y < x {
		x = y
		y = (x + n%x) % 2
	}
	return x
}
