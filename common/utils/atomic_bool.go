package utils

// https://github.com/tevino/abool

// Package abool provides atomic Boolean type for cleaner code and
// better performance.
import "sync/atomic"

// New creates an AtomicBool with default to false
func NewAtomicBool() *AtomicBool {
	return new(AtomicBool)
}

// NewBool creates an AtomicBool with given default value
func NewBool(ok bool) *AtomicBool {
	ab := NewAtomicBool()
	if ok {
		ab.Set()
	}
	return ab
}

// AtomicBool is an atomic Boolean
// Its methods are all atomic, thus safe to be called by
// multiple goroutines simultaneously
// Note: When embedding into a struct, one should always use
// *AtomicBool to avoid copy
type AtomicBool int32

// Set sets the Boolean to true
func (ab *AtomicBool) Set() {
	atomic.StoreInt32((*int32)(ab), 1)
}

// UnSet sets the Boolean to false
func (ab *AtomicBool) UnSet() {
	atomic.StoreInt32((*int32)(ab), 0)
}

// IsSet returns whether the Boolean is true
func (ab *AtomicBool) IsSet() bool {
	return atomic.LoadInt32((*int32)(ab)) == 1
}

// SetTo sets the boolean with given Boolean
func (ab *AtomicBool) SetTo(yes bool) {
	if yes {
		atomic.StoreInt32((*int32)(ab), 1)
	} else {
		atomic.StoreInt32((*int32)(ab), 0)
	}
}

// SetToIf sets the Boolean to new only if the Boolean matches the old
// Returns whether the set was done
func (ab *AtomicBool) SetToIf(old, new bool) (set bool) {
	var o, n int32
	if old {
		o = 1
	}
	if new {
		n = 1
	}
	return atomic.CompareAndSwapInt32((*int32)(ab), o, n)
}
