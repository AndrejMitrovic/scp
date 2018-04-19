package scp

import "sort"

// Value is the type of values being voted on by the network.
type Value interface {
	Less(Value) bool
	Combine(Value) Value
	Bytes() []byte
}

func VEqual(a, b Value) bool {
	return !a.Less(b) && !b.Less(a)
}

// ValueSet is a set of values, implemented as a sorted slice.
type ValueSet []Value

// Add adds a Value to a ValueSet.
// TODO: this can be done in better than O(n log n).
func (vs *ValueSet) Add(v Value) {
	if vs.Contains(v) {
		return
	}
	*vs = append(*vs, v)
	sort.Slice(*vs, func(i, j int) bool {
		return (*vs)[i].Less((*vs)[j])
	})
}

// Contains uses binary search to test whether vs contains v.
func (vs ValueSet) Contains(v Value) bool {
	if len(vs) == 0 {
		return false
	}
	mid := len(vs) / 2
	if vs[mid].Less(v) {
		if mid == len(vs)+1 {
			return false
		}
		return vs[mid+1:].Contains(v)
	}
	if v.Less(vs[mid]) {
		if mid == 0 {
			return false
		}
		return vs[:mid-1].Contains(v)
	}
	return true
}

// Combine reduces the members of vs to a single value using
// Value.Combine. The result is nil if vs is empty.
func (vs ValueSet) Combine() Value {
	if len(vs) == 0 {
		return nil
	}
	result := vs[0]
	for _, v := range vs[1:] {
		result = result.Combine(v)
	}
	return result
}
