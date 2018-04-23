package scp

// This file contains functions for finding "blocking sets" and
// "quorums" that satisfy a given predicate.
//
// Each node specifies one or more "quorum slices." Each quorum slice
// is a set of trusted peer nodes. Each quorum slice conceptually
// includes the node itself, though in this implementation that is not
// explicit.
//
// A quorum slice is not necessarily a quorum in itself. A peer in a
// quorum slice may have a dependency on a third-party node, as may
// that node, and so on. A quorum (with respect to a given node) is
// thus the transitive closure over any of its quorum slices. A node
// may have many different quorums, and they may overlap one another.
//
// Every protocol message includes the sending node's set of quorum
// slices. Every node saves the latest message seen from a given
// node. If enough messages have been seen, it is possible for a node
// to know the complete membership of one or more quorums.
//
// A "blocking set" is related to the idea of a quorum, but is
// simpler. It's any set of peers among a node's quorum slices that
// blocks the possibility of a quorum. A blocking set satisfying
// statement X precludes the existence of any quorum satisfying !X.  A
// single peer from each of a node's quorum slices is sufficient to
// form a blocking set.

func (s *Slot) findBlockingSetOrQuorum(pred predicate) NodeSet {
	nodeIDs := s.findBlockingSet(pred)
	if len(nodeIDs) > 0 {
		return nodeIDs
	}
	return s.findQuorum(pred)
}

// Checks that at least one node in each quorum slice satisfies pred.
func (s *Slot) findBlockingSet(pred predicate) NodeSet {
	var result NodeSet
	for _, slice := range s.V.Q {
		var found bool
		for _, nodeID := range slice {
			if msg, ok := s.M[nodeID]; ok && pred.test(msg) {
				found = true
				result = result.Add(nodeID)
				pred = pred.next()
				break
			}
		}
		if !found {
			return nil
		}
	}
	return result
}

// Finds a quorum in which every node satisfies the given predicate.
func (s *Slot) findQuorum(pred predicate) NodeSet {
	m := make(map[NodeID]struct{})
	m[s.V.ID] = struct{}{}
	m, _ = s.findNodeQuorum(s.V.ID, s.V.Q, pred, m)
	if len(m) == 0 {
		return nil
	}
	var result NodeSet
	for n := range m {
		result = result.Add(n)
	}
	return result
}

// Helper function for findQuorum. It checks that the given node
// (whose set of quorum slices is also given) has at least one slice
// whose members (and the transitive closure over them) all satisfy
// the given predicate.
//
// Relies on recursion (specifically, mutual recursion with
// findSliceQuorum, below) to allow backtracking. In particular, m and
// pred evolve as passing nodes are visited but must be able to revert
// to earlier values when unwinding the stack. M is the set of nodes
// visited, pred is the latest iteration of the predicate.
//
// Returns the new m and pred on success, nil and the original pred on
// failure.
func (s *Slot) findNodeQuorum(nodeID NodeID, q []NodeSet, pred predicate, m map[NodeID]struct{}) (map[NodeID]struct{}, predicate) {
	for _, slice := range q {
		m2, nextPred := s.findSliceQuorum(slice, pred, m)
		if len(m2) > 0 {
			return m2, nextPred
		}
	}
	return nil, pred
}

// Helper function for findNodeQuorum. It checks whether every node in
// a given quorum slice (and the transitive closure over them)
// satisfies the given predicate.
//
// Relies on recursion (specifically, mutual recursion with
// findNodeQuorum) to allow backtracking.
//
// Returns an updated m and pred on success, nil and the original pred
// on failure.
func (s *Slot) findSliceQuorum(slice NodeSet, pred predicate, m map[NodeID]struct{}) (map[NodeID]struct{}, predicate) {
	var newNodeIDs NodeSet // nodes in slice not yet visited (according to m)
	for _, nodeID := range slice {
		if _, ok := m[nodeID]; !ok {
			newNodeIDs = newNodeIDs.Add(nodeID)
		}
	}
	if len(newNodeIDs) == 0 {
		return m, pred
	}
	origPred := pred
	for _, nodeID := range newNodeIDs {
		if msg, ok := s.M[nodeID]; !ok || !pred.test(msg) {
			return nil, origPred
		}
		pred = pred.next()
	}
	m2 := make(map[NodeID]struct{})
	for nodeID := range m {
		m2[nodeID] = struct{}{}
	}
	for _, nodeID := range newNodeIDs {
		m2[nodeID] = struct{}{}
	}
	for _, nodeID := range newNodeIDs {
		msg := s.M[nodeID]
		m2, pred = s.findNodeQuorum(nodeID, msg.Q, pred, m2)
		if len(m2) == 0 {
			return nil, origPred
		}
	}
	return m2, pred
}

// Abstract predicate. Concrete types below.
type predicate interface {
	test(*Msg) bool

	// Allows a predicate to update itself after each successful call to
	// test, by returning a modified copy of itself for the next call to
	// test. When findQuorum needs to backtrack, it also unwinds to
	// earlier values of the predicate.
	next() predicate
}

// This is a simple function predicate. It does not change from one
// call to the next.
type fpred func(*Msg) bool

func (f fpred) test(msg *Msg) bool {
	return f(msg)
}

func (f fpred) next() predicate {
	return f
}

// This is a predicate that can narrow a set of min/max bounds as it
// traverses nodes.
type minMaxPred struct {
	min, max           int  // the current min/max bounds
	nextMin, nextMax   int  // min/max bounds for when the next predicate is generated
	finalMin, finalMax *int // each call to next updates the min/max bounds these point to
	testfn             func(msg *Msg, min, max int) (bool, int, int)
}

func (p *minMaxPred) test(msg *Msg) bool {
	p.nextMin, p.nextMax = p.min, p.max
	if p.min > p.max {
		return false
	}
	res, min, max := p.testfn(msg, p.min, p.max)
	if !res {
		return false
	}
	p.nextMin, p.nextMax = min, max
	return true
}

func (p *minMaxPred) next() predicate {
	if p.finalMin != nil {
		*p.finalMin = p.nextMin
	}
	if p.finalMax != nil {
		*p.finalMax = p.nextMax
	}
	return &minMaxPred{
		min:      p.nextMin,
		max:      p.nextMax,
		finalMin: p.finalMin,
		finalMax: p.finalMax,
		testfn:   p.testfn,
	}
}
