package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/bobg/scp"

	"github.com/chain/txvm/crypto/ed25519"
)

type (
	marshaled struct {
		M json.RawMessage
		S string // hex-encoded signature over marshaledPayload
	}

	marshaledPayload struct {
		C int32
		V string
		I int
		Q [][]string
		T marshaledTopic
	}

	marshaledTopic struct {
		Type        int // scp.Phase values
		X, Y        []string
		B, C, P, PP marshaledBallot
		PN, HN, CN  int
	}

	marshaledBallot struct {
		N int
		X string
	}
)

func marshal(msg *scp.Msg) ([]byte, error) {
	var q [][]string
	for _, slice := range msg.Q {
		var qslice []string
		for _, id := range slice {
			qslice = append(qslice, string(id))
		}
		q = append(q, qslice)
	}

	var mt marshaledTopic
	switch topic := msg.T.(type) {
	case *scp.NomTopic:
		var x, y []string
		for _, val := range topic.X {
			x = append(x, val)
		}
		for _, val := range topic.Y {
			y = append(y, val)
		}
		mt.X = x
		mt.Y = y

	case *scp.PrepTopic:
		mt.B = marshaledBallot{N: topic.B.N, X: topic.B.X.(valtype).String()}
		mt.P = marshaledBallot{N: topic.P.N, X: topic.P.X.(valtype).String()}
		mt.PP = marshaledBallot{N: topic.PP.N, X: topic.PP.X.(valtype).String()}
		mt.HN = topic.HN
		mt.CN = topic.CN

	case *scp.CommitTopic:
		mt.B = marshaledBallot{N: topic.B.N, X: topic.B.X.(valtype).String()}
		mt.PN = topic.PN
		mt.HN = topic.HN
		mt.CN = topic.CN

	case *scp.ExtTopic:
		mt.C = marshaledBallot{N: topic.C.N, X: topic.C.X.(valtype).String()}
		mt.HN = topic.HN
	}
	mp := marshaledPayload{
		C: msg.C,
		V: string(msg.V),
		I: int(msg.I),
		Q: q,
		T: mt,
	}
	mpbytes, err := json.Marshal(mp) // xxx json is subject to mutation in transit!
	if err != nil {
		return nil, err
	}
	sig := ed25519.Sign(prv, mpbytes)
	m := marshaled{
		M: mpbytes,
		S: hex.EncodeToString(sig),
	}
	return json.Marshal(m)
}

func unmarshal(b []byte) (*scp.Msg, error) {
	var m marshaled
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	var mp marshaledPayload
	err = json.Unmarshal(m.M, &mp)
	if err != nil {
		return nil, err
	}

	sig, err := hex.DecodeString(m.S)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(mp.V)
	if err != nil {
		return nil, err
	}
	pubkeyHex := u.Path
	strings.Trim(pubkeyHex, "/")
	pubkey, err := hex.DecodeString(pubkeyHex)
	if err != nil {
		return nil, err
	}
	if !ed25519.Verify(pubkey, m.M, sig) {
		return nil, errors.New("bad signature")
	}

	var q []scp.NodeIDSet
	for _, slice := range mp.Q {
		var qslice scp.NodeIDSet
		for _, id := range slice {
			qslice = qslice.Add(scp.NodeID(id))
		}
		q = append(q, qslice)
	}

	var topic scp.Topic
	switch scp.Phase(mp.T.Type) {
	case scp.PhNom:
		topic = &scp.NomTopic{
			X: x,
			Y: y,
		}

	case scp.PhPrep:
		topic = &scp.PrepTopic{
			B:  b,
			P:  p,
			PP: pp,
			HN: mp.T.HN,
			CN: mp.T.CN,
		}

	case scp.PhCommit:
		topic = &scp.CommitTopic{
			B:  b,
			PN: mp.T.PN,
			HN: mp.T.HN,
			CN: mp.T.CN,
		}

	case scp.PhExt:
		topic = &scp.ExtTopic{
			C:  c,
			HN: mp.T.HN,
		}

	default:
		return nil, fmt.Errorf("unknown topic type %d", mp.T.Type)
	}

	msg := &scp.Msg{
		C: mp.C,
		V: mp.V,
		I: mp.I,
		Q: q,
		T: topic,
	}
	return msg, nil
}
