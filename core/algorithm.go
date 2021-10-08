// ---------------------------------------------------------------------------
// LN-ELECTRONIC PROJECT LN-16C10R
// wwww.ln-electronic.com  ProjectManager : @Kimera
// FileData: 8/10/2564 7:05 2564  FileName : algorithm.go
// ---------------------------------------------------------------------------

package core

import (
	"errors"
	"time"
)

// These constants are the bit lengths of lnalgorithm ID parts.
const (
	TimestampLength = 41
	MachineIDLength = 10
	SequenceLength  = 12
	MaxSequence     = 1<<SequenceLength - 1
	MaxTimestamp    = 1<<TimestampLength - 1
	MaxMachineID    = 1<<MachineIDLength - 1

	machineIDMoveLength = SequenceLength
	timestampMoveLength = MachineIDLength + SequenceLength
)

//   AtomicResolver : base sync/atomic (by default).
type SequenceResolver func(ms int64) (uint16, error)

// default machineID is 0
// default resolver is AtomicResolver
var (
	resolver  SequenceResolver
	machineID = 0
	startTime = time.Date(2008, 11, 10, 23, 0, 0, 0, time.UTC)
)

func ID() uint64 {
	id, _ := NextID()
	return id
}

func NextID() (uint64, error) {
	c := currentMillis()
	seqResolver := callSequenceResolver()
	seq, err := seqResolver(c)

	if err != nil {
		return 0, err
	}

	for seq >= MaxSequence {
		c = waitForNextMillis(c)
		seq, err = seqResolver(c)
		if err != nil {
			return 0, err
		}
	}

	df := int(elapsedTime(c, startTime))
	if df < 0 || df > MaxTimestamp {
		return 0, errors.New("The maximum life cycle of the Ln-algorithm is 2^41-1(millis), please check starttime")
	}

	id := uint64((df << timestampMoveLength) | (machineID << machineIDMoveLength) | int(seq))
	return id, nil
}

// SetStartTime set the start time for lnalgorithm algorithm.

func SetStartTime(s time.Time) {
	s = s.UTC()

	if s.IsZero() {
		panic("The start time cannot be a zero value")
	}

	if s.After(time.Now()) {
		panic("The s cannot be greater than the current millisecond")
	}

	// Because s must after now, so the `df` not < 0.
	df := elapsedTime(currentMillis(), s)
	if df > MaxTimestamp {
		panic("The maximum life cycle of the ln algorithm is 69 years")
	}

	startTime = s
}

func SetMachineID(m uint16) {
	if m > MaxMachineID {
		panic("The machineid cannot be greater than 1023")
	}
	machineID = int(m)
}

// SetSequenceResolver set an custom sequence resolver.
func SetSequenceResolver(seq SequenceResolver) {
	if seq != nil {
		resolver = seq
	}
}

// SID lnalgorithm id
type SID struct {
	Sequence  uint64
	MachineID uint64
	Timestamp uint64
	ID        uint64
}

// GenerateTime lnalgorithm generate at, return a UTC time.
func (id *SID) GenerateTime() time.Time {
	ms := startTime.UTC().UnixNano()/1e6 + int64(id.Timestamp)

	return time.Unix(0, (ms * int64(time.Millisecond))).UTC()
}

// ParseID parse lnalgorithm it to SID struct.
func ParseID(id uint64) SID {
	time := id >> (SequenceLength + MachineIDLength)
	sequence := id & MaxSequence
	machineID := (id & (MaxMachineID << SequenceLength)) >> SequenceLength

	return SID{
		ID:        id,
		Sequence:  sequence,
		MachineID: machineID,
		Timestamp: time,
	}
}

//--------------------------------------------------------------------
// private function defined.
//--------------------------------------------------------------------

func waitForNextMillis(last int64) int64 {
	now := currentMillis()
	for now == last {
		now = currentMillis()
	}
	return now
}

func callSequenceResolver() SequenceResolver {
	if resolver == nil {
		return AtomicResolver
	}

	return resolver
}

func elapsedTime(nowms int64, s time.Time) int64 {
	return nowms - s.UTC().UnixNano()/1e6
}

// currentMillis get current millisecond.
func currentMillis() int64 {
	return time.Now().UTC().UnixNano() / 1e6
}
