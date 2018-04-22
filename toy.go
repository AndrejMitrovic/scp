// +build ignore

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"log"
	"math/rand"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bobg/scp"
)

type entry struct {
	node *scp.Node
	ch   chan *scp.Env
}

type valType string

func (v valType) Less(other scp.Value) bool {
	return v < other.(valType)
}

func (v valType) Combine(other scp.Value) scp.Value {
	if v < other.(valType) {
		return v
	}
	return other
}

func (v valType) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, v)
	return buf.Bytes()
}

func (v valType) String() string {
	return string(v)
}

// Usage:
//   go run toy.go [-seed N] 'alice: bob carol david / bob carol ed / fran gabe hank' 'bob: alice carol david / gabe hank' ...
// Each argument gives a node's name (before the colon) and the node's
// quorum slices.
// Nodes do not specify themselves as quorum slice members, though
// they are understood to belong to every quorum slice.

func main() {
	seed := flag.Int64("seed", 1, "RNG seed")
	flag.Parse()
	rand.Seed(*seed)

	entries := make(map[scp.NodeID]entry)

	ch := make(chan *scp.Env, 10000)
	var highestSlot int32
	for _, arg := range flag.Args() {
		parts := strings.SplitN(arg, ":", 2)
		nodeID := scp.NodeID(parts[0])
		var q [][]scp.NodeID
		for _, slices := range strings.Split(parts[1], "/") {
			var qslice []scp.NodeID
			for _, field := range strings.Fields(slices) {
				if field == string(nodeID) {
					log.Print("skipping quorum slice member %s for node %s", field, nodeID)
					continue
				}
				qslice = append(qslice, scp.NodeID(field))
			}
			q = append(q, qslice)
		}
		node := scp.NewNode(nodeID, q)
		nodeCh := make(chan *scp.Env, 1000)
		entries[nodeID] = entry{node: node, ch: nodeCh}
		go nodefn(node, nodeCh, ch, &highestSlot)
	}

	// log.Print(spew.Sdump(entries))

	for env := range ch {
		if _, ok := env.M.(*scp.NomMsg); !ok && int32(env.I) > highestSlot { // this is the only thread that writes highestSlot, so it's ok to read it non-atomically
			atomic.StoreInt32(&highestSlot, int32(env.I))
			log.Printf("highestSlot is now %d", highestSlot)
		}

		// Send this message to every other node.
		for nodeID, entry := range entries {
			if nodeID == env.V {
				continue
			}
			entry.ch <- env
		}
	}
}

const (
	minNomDelayMS = 500
	maxNomDelayMS = 2000
)

// runs as a goroutine
func nodefn(n *scp.Node, recv <-chan *scp.Env, send chan<- *scp.Env, highestSlot *int32) {
	for {
		// Some time in the next minNomDelayMS to maxNomDelayMS
		// milliseconds.
		timer := time.NewTimer(time.Duration((minNomDelayMS + rand.Intn(maxNomDelayMS-minNomDelayMS)) * int(time.Millisecond)))

		select {
		case env := <-recv:
			// Never mind about the nomination timer.
			if !timer.Stop() {
				<-timer.C
			}

			res, err := n.Handle(env)
			if err != nil {
				n.Logf("could not handle %s: %s", env, err)
				continue
			}
			if res != nil {
				n.Logf("handled %s -> %s", env, res)
				send <- res
			}

		case <-timer.C:
			val := foods[rand.Intn(len(foods))]
			slotID := scp.SlotID(atomic.LoadInt32(highestSlot))

			if slot, ok := n.Pending[slotID]; ok && slot.Ph == scp.PhNom {
				// Prod the node to try to get it out of NOMINATE phase, by
				// resending it all the messages it's received so far. If it's
				// a later round than the last time the node was prodded, the
				// outcome may be different.
				slot.Logf("prodding")
				for _, env := range slot.M {
					res, err := n.Handle(env)
					if err != nil {
						n.Logf("error prodding node with %s: %s", env, err)
					} else if res != nil {
						send <- res
					}
				}
				continue
			}

			// Wanna be startin' something.
			slotID++

			// Send a nominate message "from" the node to itself. If it has
			// max priority among its neighbors (for this slot) it will
			// propagate the nomination.
			var vs scp.ValueSet
			vs = vs.Add(val)
			env := scp.NewEnv(n.ID, scp.SlotID(slotID), n.Q, &scp.NomMsg{X: vs})
			n.Logf("trying to get something started with %s", env)
			res, err := n.Handle(env)
			if err != nil {
				n.Logf("could not handle %s: %s", env, err)
				continue
			}
			if res != nil {
				send <- res
			}
		}
	}
}

var foods = []valType{
	"pizza",
	"burgers",
	"burritos",
	"sandwiches",
	"sushi",
	"salads",
	"gyros",
	"indian",
	"soup",
	"pasta",
}
