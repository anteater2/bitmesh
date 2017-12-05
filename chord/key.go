package chord

import (
	"hash/fnv"
)

// Key is a key in the distributed hash table.
// IT MUST BE BOUNDED BY MAXKEY
type Key uint64

// BetweenExclusive returns if a key is in (start, end)
// Note that it is possible for the interval to start and end at the same key
// The interval is just the clockwise sweep between start and end.
func (key Key) BetweenExclusive(start Key, end Key) bool {
	s, e := uint64(start), uint64(end)
	k := uint64(key)
	if s > config.maxKey || e > config.maxKey || k > config.maxKey {
		panic("MaxKey constraint has been violated!")
	}
	if s == e {
		return k != s && k != e // Full sweep - all keys are in range, unless it is s or e.
	} else if s > e { // Interval wraps - if key is lt end or gt start, it is in interval
		return s < k || k < e
	} else {
		return s < k && k < e
	}
}

// BetweenEndInclusive returns if a key is in (start,end]
// Note that it is possible for the interval to start and end at the same key
// The interval is just the clockwise sweep between start and end.
func (key Key) BetweenEndInclusive(start Key, end Key) bool {
	s, e := uint64(start), uint64(end)
	k := uint64(key)
	if s > config.maxKey || e > config.maxKey || k > config.maxKey {
		panic("MaxKey constraint has been violated!")
	}
	if s == e {
		return true // Full sweep - all keys are in range.
	}
	if s > e { // Interval wraps - if key is lt end or gt start, it is in interval
		return s < k || k <= e
	} else {
		return (s < k && k <= e)
	}
}

// NewKey returns a new key
func NewKey(value uint64) Key {
	return Key(value)
}

// Valid returns true if the key is within the keyspace, false otherwise
func (key Key) Valid() bool {
	k := uint64(key)
	return 0 <= k && k < config.maxKey // Check for 0 <= not necessary
}

// Hash a string, returning a key bounded by maxKey.
func Hash(s string, maxKey uint64) Key {
	h := fnv.New64a()
	h.Write([]byte(s))
	return NewKey(h.Sum64() % maxKey)
}
