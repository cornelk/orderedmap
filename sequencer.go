package orderedmap

import "sync/atomic"

// sequence is a global sequencer that keeps track of any call to an Unmarshal function.
// This can not be part of the map type itself as the recursive unmarshalling will call
// the Unmarshal function of a new Entry type, which does not have access to the map it is
// part of.
var sequence uint64

func nextSequence() uint64 {
	return atomic.AddUint64(&sequence, 1)
}
