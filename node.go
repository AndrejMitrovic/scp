package scp

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"math/big"
	"sync"

	"github.com/davecgh/go-xdr/xdr"
)

type NodeID interface {
	String() string
}

type Node struct {
	ID      NodeID
	Q       *QSet
	Pending map[SlotID]*Slot
	Ext     map[SlotID]*ExtMsg

	mu sync.Mutex
}

func NewNode(id NodeID, q *QSet) *Node {
	return &Node{
		ID:      id,
		Q:       q,
		Pending: make(map[SlotID]*Slot),
		Ext:     make(map[SlotID]*ExtMsg),
	}
}

func (n *Node) Handle(env *Env, ch chan<- *Env) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if msg, ok := n.Ext[env.I]; ok {
		// This node has already externalized a value for the given slot.
		// Send an EXTERNALIZE message outbound, unless the inbound
		// message is also EXTERNALIZE.
		// TODO: ...in which case double-check that the values agree?
		if _, ok := env.M.(*ExtMsg); !ok {
			ch <- &Env{
				V: n.ID,
				I: env.I,
				Q: n.Q,
				M: msg,
			}
		}
		return
	}

	s, ok := n.Pending[env.I]
	if !ok {
		s = &Slot{ID: env.I, V: n}
		n.Pending[env.I] = s
	}

	outbound := s.Handle(env)
	if outbound == nil {
		return
	}

	if extMsg, ok := outbound.M.(*ExtMsg); ok {
		// Handling the inbound message resulted in externalizing a value.
		// We can now save the EXTERNALIZE message and get rid of the Slot
		// object.
		n.Ext[env.I] = extMsg
		delete(n.Pending, env.I)
	}

	ch <- outbound
}

var ErrNoPrev = errors.New("no previous value")

func (n *Node) G(i SlotID, m []byte) (result [32]byte, err error) {
	hasher := sha256.New()

	var prevValBytes []byte
	if i > 1 {
		msg, ok := n.Ext[i]
		if !ok {
			return result, ErrNoPrev
		}
		prevValBytes = msg.C.X.Bytes()
	}

	r, _ := xdr.Marshal(i)
	hasher.Write(r)
	hasher.Write(prevValBytes)
	hasher.Write(m)
	hasher.Sum(result[:0])

	return result, nil
}

func (n *Node) Weight(id NodeID) float64 {
	if id == n.ID {
		return 1.0
	}
	count := 0
	n.Q.Each(func(ids []NodeID) {
		for _, thisID := range ids {
			if id == thisID {
				count++
				break
			}
		}
	})
	return float64(count) / float64(n.Q.Size())
}

// maxuint256, as a float
var hmax *big.Float

func (n *Node) Neighbors(i SlotID, num int) ([]NodeID, error) {
	peers := n.Peers()
	peers = append(peers, node.ID)
	var result []NodeID
	for _, nodeID := range peers {
		w := big.NewFloat(n.Weight(nodeID))
		w.Mul(w, hmax)
		hwInt, _ := w.Int(nil)
		hwBytes := hwInt.Bytes()
		var hw [32]byte
		copy(hw[32-len(hwBytes):], hwBytes) // hw is now a big-endian uint256

		m := new(bytes.Buffer)
		m.WriteByte('N')
		numBytes, _ := xdr.Marshal(num)
		m.Write(numBytes)
		m.WriteString(nodeID)
		g, err := n.G(m.Bytes())
		if err != nil {
			return nil, err
		}
		if bytes.Compare(g[:], hw[:]) < 0 {
			result = append(result, nodeID)
		}
	}
	return result, nil
}

func init() {
	maxUint256 := [32]byte{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}
	hmaxInt := new(big.Int)
	hmaxInt.SetBytes(maxUint256[:])
	hmax = new(big.Float)
	hmax.SetInt(hmaxInt)
}
