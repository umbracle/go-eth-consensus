package spec

import (
	"path/filepath"
	"reflect"
	"testing"

	consensus "github.com/umbracle/go-eth-consensus"
)

type epochProcessignFunc func(state *consensus.BeaconStatePhase0) error

func TestEpochProcessing(t *testing.T) {
	type epochTest struct {
		Pre  consensus.BeaconStatePhase0
		Post consensus.BeaconStatePhase0
	}

	cases := []struct {
		name    string
		path    string
		handler epochProcessignFunc
	}{
		{
			"Effective balance updates",
			"effective_balance_updates/*/*",
			processEffectiveBalanceUpdates,
		},
		{
			"Eth1 data reset",
			"eth1_data_reset/*/*",
			processEth1DataReset,
		},
		{
			"Historical roots update",
			"historical_roots_update/*/*",
			processHistoricalRootsUpdate,
		},
		{
			"Justification_and_finalization",
			"justification_and_finalization/*/*",
			processJustificationAndFinalization,
		},
		{
			"Participation record",
			"participation_record_updates/*/*",
			processParticipationRecordUpdates,
		},
		{
			"Randao mix",
			"randao_mixes_reset/*/*",
			processRandaoMixesReset,
		},
		{
			"Registry updates",
			"registry_updates/*/*",
			processRegistryUpdates,
		},
		{
			"Rewards and Penalties",
			"rewards_and_penalties/*/*",
			processRewardsAndPenalties,
		},
		{
			"Process slashing",
			"slashings/*/*",
			processSlashings,
		},
		{
			"Slashings reset",
			"slashings_reset/*/*",
			processSlashingsReset,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			listTestData(t, filepath.Join("mainnet/phase0/epoch_processing/", c.path), func(th *testHandler) {
				eTest := &epochTest{}
				th.decodeFile("pre", &eTest.Pre)
				ok := th.decodeFile("post", &eTest.Post, true)

				if err := c.handler(&eTest.Pre); err != nil {
					if ok {
						t.Fatal(err)
					}
					return
				}

				if !ok {
					t.Fatal("it should fail")
				}
				if !reflect.DeepEqual(eTest.Pre, eTest.Post) {
					t.Fatal("bad")
				}
			})
		})
	}
}
