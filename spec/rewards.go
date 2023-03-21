package spec

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"

	ssz "github.com/ferranbt/fastssz"
	eth2_shuffle "github.com/protolambda/eth2-shuffle"
	consensus "github.com/umbracle/go-eth-consensus"
	"github.com/umbracle/go-eth-consensus/bitlist"
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

func getMatchingHeadAttestations(state *consensus.BeaconStatePhase0, epoch uint64) []*consensus.PendingAttestation {
	res := []*consensus.PendingAttestation{}
	for _, a := range getMatchingTargetAttestations(state, epoch) {
		root := getBlockRootAtSlot(state, a.Data.Slot)

		if bytes.Equal(a.Data.BeaconBlockHash[:], root[:]) {
			res = append(res, a)
		}
	}

	return res
}

func getMatchingSourceAttestations(state *consensus.BeaconStatePhase0, epoch uint64) []*consensus.PendingAttestation {
	if epoch == getCurrentEpoch(state) {
		return state.CurrentEpochAttestations
	}
	return state.PreviousEpochAttestations
}

func getTargetDeltas(state *consensus.BeaconStatePhase0) ([]uint64, []uint64) {
	return getAttestationComponentDeltas(state, getMatchingTargetAttestations(state, getPreviousEpoch(state)))
}

func getHeadDeltas(state *consensus.BeaconStatePhase0) ([]uint64, []uint64) {
	return getAttestationComponentDeltas(state, getMatchingHeadAttestations(state, getPreviousEpoch(state)))
}

func getSourceDeltas(state *consensus.BeaconStatePhase0) ([]uint64, []uint64) {
	return getAttestationComponentDeltas(state, getMatchingSourceAttestations(state, getPreviousEpoch(state)))
}

func getPreviousEpoch(state *consensus.BeaconStatePhase0) uint64 {
	curEpoch := getCurrentEpoch(state)
	if curEpoch == 0 {
		return 0
	}
	return curEpoch - 1
}

func getCurrentEpoch(state *consensus.BeaconStatePhase0) uint64 {
	return state.Slot / Spec.SlotsPerEpoch
}

func getAttestationComponentDeltas(state *consensus.BeaconStatePhase0, attestations []*consensus.PendingAttestation) ([]uint64, []uint64) {
	numValidators := len(state.Validators)
	rewards := make([]uint64, numValidators)
	penalties := make([]uint64, numValidators)

	totalBalance := getTotalActiveBalance(state)

	unslashedAttestingIndices, err := getUnslashedAttestingIndices(state, attestations)
	if err != nil {
		panic(err)
	}

	attestingBalance := getTotalBalance(state, unslashedAttestingIndices)

	unslashedAttestingIndicesMap := make(map[uint64]bool)
	for _, i := range unslashedAttestingIndices {
		unslashedAttestingIndicesMap[i] = true
	}

	for _, indx := range getElegibleValidatorIndices(state) {
		if unslashedAttestingIndicesMap[indx] {
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

	return effectiveBalance * Spec.BaseRewardFactor / integerSquareRoot(totalBalance) / Spec.BaseRewardsPerEpoch
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

func getUnslashedAttestingIndices(state *consensus.BeaconStatePhase0, attestations []*consensus.PendingAttestation) ([]uint64, error) {
	output := make([]uint64, 0)
	seen := make(map[uint64]bool)

	for _, a := range attestations {
		indices, err := getAttestingIndices(state, a.Data, a.AggregationBits)
		if err != nil {
			return nil, err
		}
		for _, i := range indices {
			if !seen[i] {
				output = append(output, i)
			}
			seen[i] = true
		}
	}

	// Sort the attesting set indices by increasing order.
	sort.Slice(output, func(i, j int) bool {
		return output[i] < output[j]
	})

	// Remove slashed validator indices.
	ret := make([]uint64, 0)
	for i := range output {
		val := state.Validators[output[i]]
		if !val.Slashed {
			ret = append(ret, output[i])
		}
	}
	return ret, nil
}

func getAttestingIndices(state *consensus.BeaconStatePhase0, data *consensus.AttestationData, bits []byte) ([]uint64, error) {
	blist := bitlist.BitList(bits)

	committee := getBeaconCommittee(state, data.Slot, data.Index)

	if blist.Len() != uint64(len(committee)) {
		return nil, fmt.Errorf("bad size")
	}

	res := []uint64{}
	for indx, c := range committee {
		if blist.BitAt(uint64(indx)) {
			res = append(res, c)
		}
	}
	return res, nil
}

func getBeaconCommittee(state *consensus.BeaconStatePhase0, slot uint64, index uint64) []uint64 {
	epoch := computeEpochAtSlot(slot)
	committeesPerSlot := getCommitteeCountPerSlot(state, epoch)

	seed := getSeed(state, epoch, consensus.DomainBeaconAttesterType)
	active := getActiveValidatorIndices(state, epoch)

	return computeCommittee(
		active,
		seed,
		(slot%Spec.SlotsPerEpoch)*committeesPerSlot+index,
		committeesPerSlot*Spec.SlotsPerEpoch,
	)
}

func getCommitteeCountPerSlot(state *consensus.BeaconStatePhase0, epoch uint64) uint64 {
	return max(1, min(Spec.MaxCommitteesPerSlot, uint64(len(getActiveValidatorIndices(state, epoch)))/Spec.SlotsPerEpoch/Spec.TargetCommitteeSize))
}

func getSeed(state *consensus.BeaconStatePhase0, epoch uint64, domain consensus.Domain) consensus.Root {
	mix := getRandaoMix(state, epoch+Spec.EpochsPerHistoricalVector-Spec.MinSeedLookAhead-1)

	epochBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(epochBuf, epoch)

	hash := sha256.New()
	hash.Write(domain[:])
	hash.Write(epochBuf)
	hash.Write(mix[:])

	root := consensus.Root{}
	copy(root[:], hash.Sum(nil))

	return root
}

func getRandaoMix(state *consensus.BeaconStatePhase0, epoch uint64) [32]byte {
	return state.RandaoMixes[epoch%Spec.EpochsPerHistoricalVector]
}

func max(i, j uint64) uint64 {
	if i > j {
		return i
	}
	return j
}

func min(i, j uint64) uint64 {
	if i < j {
		return i
	}
	return j
}

func computeEpochAtSlot(slot uint64) uint64 {
	return slot / Spec.SlotsPerEpoch
}

func getActiveValidatorIndices(state *consensus.BeaconStatePhase0, epoch uint64) []uint64 {
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
	return getTotalBalance(state, getActiveValidatorIndices(state, getCurrentEpoch(state)))
}

func getTotalBalance(state *consensus.BeaconStatePhase0, indices []uint64) uint64 {
	balance := uint64(0)

	for _, indx := range indices {
		balance += state.Validators[indx].EffectiveBalance
	}

	balance = max(balance, Spec.EffectiveBalanceIncrement)
	return balance
}

func computeCommittee(indices []uint64, seed consensus.Root, index, count uint64) []uint64 {
	numActiveValidators := uint64(len(indices))

	start := (numActiveValidators * index) / count
	end := (numActiveValidators * (index + 1)) / count

	eth2ShuffleHashFunc := func(data []byte) []byte {
		hash := sha256.New()
		hash.Write(data)
		return hash.Sum(nil)
	}

	commmittee := make([]uint64, len(indices))
	copy(commmittee[:], indices)

	eth2_shuffle.UnshuffleList(eth2ShuffleHashFunc, commmittee, shuffleRoundCount, seed)

	return commmittee[start:end]
}

func integerSquareRoot(n uint64) uint64 {
	x := n
	y := (x + 1) >> 1

	for y < x {
		x = y
		y = (x + n/x) >> 1
	}
	return x
}
