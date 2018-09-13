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

// Checks that at least one node in each quorum slice satisfies pred
// (excluding the slot's node).
func (s *Slot) findBlockingSet(pred predicate) NodeIDSet {
	res, _ := s.V.Q.findBlockingSet(s.M, pred)
	return res
}

// Finds a quorum in which every node satisfies the given
// predicate. The slot's node itself is presumed to satisfy the
// predicate.
func (s *Slot) findQuorum(pred predicate) NodeIDSet {
	res, _ := s.V.Q.findQuorum(s.V.ID, s.M, pred)
	return res
}

// Tells whether a statement can be accepted, either because a
// blocking set accepts, or because a quorum votes-or-accepts it. The
// function f should produce an "accepts" predicate when its argument
// is false and a "votes-or-accepts" predicate when its argument is
// true.
func (s *Slot) accept(f func(bool) predicate) NodeIDSet {
	// 1. If s's node accepts the statement,
	//    we're done
	//    (since it is its own blocking set and,
	//    more intuitively,
	//    node N can accept X if N already accepts X).
	if s.sent != nil && f(false).test(s.sent) {
		return NodeIDSet{s.V.ID}
	}

	// 2. Look for a blocking set apart from s.V that accepts.
	nodeIDs := s.findBlockingSet(f(false))
	if len(nodeIDs) > 0 {
		return nodeIDs
	}

	// 3. Look for a quorum that votes-or-accepts.
	return s.findQuorum(f(true))
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

// This is a predicate that can narrow a set of values as it traverses
// nodes.
type valueSetPred struct {
	vals      ValueSet
	nextVals  ValueSet
	finalVals *ValueSet
	testfn    func(*Msg, ValueSet) ValueSet
}

func (p *valueSetPred) test(msg *Msg) bool {
	if len(p.vals) == 0 {
		return false
	}
	p.nextVals = p.testfn(msg, p.vals)
	return len(p.nextVals) > 0
}

func (p *valueSetPred) next() predicate {
	if p.finalVals != nil {
		*p.finalVals = p.nextVals
	}
	return &valueSetPred{
		vals:      p.nextVals,
		finalVals: p.finalVals,
		testfn:    p.testfn,
	}
}

// This is a predicate that can narrow a set of ballots as it traverses
// nodes.
type ballotSetPred struct {
	ballots      BallotSet
	nextBallots  BallotSet
	finalBallots *BallotSet
	testfn       func(*Msg, BallotSet) BallotSet
}

func (p *ballotSetPred) test(msg *Msg) bool {
	if len(p.ballots) == 0 {
		return false
	}
	p.nextBallots = p.testfn(msg, p.ballots)
	return len(p.nextBallots) > 0
}

func (p *ballotSetPred) next() predicate {
	if p.finalBallots != nil {
		*p.finalBallots = p.nextBallots
	}
	return &ballotSetPred{
		ballots:      p.nextBallots,
		finalBallots: p.finalBallots,
		testfn:       p.testfn,
	}
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
