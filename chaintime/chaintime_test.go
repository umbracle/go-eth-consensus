package chaintime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func restoreHooks() func() {
	n := now
	return func() {
		now = n
	}
}

func TestChainTime_SlotToEpoch(t *testing.T) {
	c := New(time.Now(), 0, 10)

	var cases = []struct {
		slot  uint64
		epoch uint64
	}{
		{0, 0},
		{1, 0},
		{9, 0},
		{10, 1},
		{11, 1},
		{19, 1},
		{20, 2},
		{21, 2},
	}
	for _, cc := range cases {
		assert.Equal(t, c.SlotToEpoch(cc.slot), cc.epoch)
	}
}

func TestChainTime_Active(t *testing.T) {
	defer restoreHooks()

	c := New(time.Unix(10, 0), 10, 1)

	now = func() time.Time { return time.Unix(0, 0) }
	assert.False(t, c.IsActive())

	now = func() time.Time { return time.Unix(11, 0) }
	assert.True(t, c.IsActive())
}

func TestChainTime_Current(t *testing.T) {
	defer restoreHooks()

	c := New(time.Unix(0, 0), 10, 1)

	now = func() time.Time { return time.Unix(21, 0) }
	e := c.CurrentEpoch()
	assert.Equal(t, e.Number, uint64(2))

	s := c.CurrentSlot()
	assert.Equal(t, s.Number, uint64(2))
}

func TestChainTime_GetEpoch(t *testing.T) {
	defer restoreHooks()

	c := New(time.Unix(10, 0), 10, 10)

	s := c.Epoch(3)
	assert.Equal(t, s.Number, uint64(3))

	expectedTime := int64(10 + 3*10*10)
	assert.Equal(t, s.Time, time.Unix(expectedTime, 0))

	// one second to epoch
	now = func() time.Time { return time.Unix(expectedTime-1, 0) }

	select {
	case <-s.C().C:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout")
	}
}

func TestChainTime_GetSlot(t *testing.T) {
	defer restoreHooks()

	c := New(time.Unix(10, 0), 10, 10)

	s := c.Slot(20)
	assert.Equal(t, s.Number, uint64(20))
	assert.Equal(t, s.Epoch, uint64(2))

	expectedTime := int64(10 + 20*10)
	assert.Equal(t, s.Time, time.Unix(expectedTime, 0))

	// one second to slot time
	now = func() time.Time { return time.Unix(expectedTime-1, 0) }

	select {
	case <-s.C().C:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout")
	}
}
