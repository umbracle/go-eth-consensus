package chaintime

import (
	"time"
)

type Chaintime struct {
	Genesis        time.Time
	SecondsPerSlot uint64
	SlotsPerEpoch  uint64
}

func New(genesis time.Time, secondsPerSlot, slotsPerEpoch uint64) *Chaintime {
	return &Chaintime{
		Genesis:        genesis,
		SecondsPerSlot: secondsPerSlot,
		SlotsPerEpoch:  slotsPerEpoch,
	}
}

func (c *Chaintime) IsActive() bool {
	return c.Genesis.Before(now())
}

func (c *Chaintime) SlotToEpoch(slot uint64) uint64 {
	return slot / c.SlotsPerEpoch
}

func (c *Chaintime) newTime(seconds uint64) time.Time {
	return c.Genesis.Add(time.Duration(seconds) * time.Second)
}

func (c *Chaintime) CurrentEpoch() Epoch {
	numEpoch := uint64(now().Sub(c.Genesis).Seconds()) / c.SecondsPerSlot * c.SlotsPerEpoch
	return c.Epoch(numEpoch)
}

func (c *Chaintime) CurrentSlot() Slot {
	numSlot := uint64(now().Sub(c.Genesis).Seconds()) / c.SecondsPerSlot
	return c.Slot(numSlot)
}

func (c *Chaintime) Slot(slot uint64) Slot {
	s := Slot{
		Number: slot,
		Time:   c.newTime(slot * c.SecondsPerSlot),
		Epoch:  c.SlotToEpoch(slot),
	}
	return s
}

func (c *Chaintime) Epoch(epoch uint64) Epoch {
	e := Epoch{
		Number: epoch,
		Time:   c.newTime(epoch * c.SlotsPerEpoch * c.SecondsPerSlot),
	}
	return e
}

type Epoch struct {
	Number uint64
	Time   time.Time
}

func (e Epoch) C() *time.Timer {
	return time.NewTimer(e.Time.Sub(now()))
}

type Slot struct {
	Number uint64
	Time   time.Time
	Epoch  uint64
}

func (s Slot) C() *time.Timer {
	return time.NewTimer(s.Time.Sub(now()))
}

var now = time.Now
