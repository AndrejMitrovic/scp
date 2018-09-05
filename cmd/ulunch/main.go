package main

// Usage:
//   ulunch [-seed N] CONFIGFILE

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"

	"github.com/BurntSushi/toml"
	"github.com/bobg/scp"
)

type valType string

func (v valType) MarshalJSON() ([]byte, error) {
	return []byte(v), nil
}

func (v *valType) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, (*string)(v))
}

func (v valType) Less(other scp.Value) bool {
	return v < other.(valType)
}

func (v valType) Combine(other scp.Value, slotID scp.SlotID) scp.Value {
	if slotID%2 == 0 {
		if v > other.(valType) {
			return v
		}
	} else if v < other.(valType) {
		return v
	}
	return other
}

func (v valType) IsNil() bool {
	return v == ""
}

func (v valType) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, v)
	return buf.Bytes()
}

func (v valType) String() string {
	return string(v)
}

func main() {
	seed := flag.Int64("seed", 1, "RNG seed")
	delay := flag.Int("delay", 100, "random delay limit in milliseconds")
	flag.Parse()
	rand.Seed(*seed)

	if flag.NArg() < 1 {
		log.Fatal("usage: lunch [-seed N] CONFFILE")
	}
	confFile := flag.Arg(0)
	confBits, err := ioutil.ReadFile(confFile)
	if err != nil {
		log.Fatal(err)
	}
	var conf struct {
		Nodes map[string][][]string
	}
	_, err = toml.Decode(string(confBits), &conf)
	if err != nil {
		log.Fatal(err)
	}

	nodes := make(map[scp.NodeID]*scp.Node)
	ch := make(chan *scp.Msg)
	for nodeID, qstrs := range conf.Nodes {
		q := make([]scp.NodeIDSet, 0, len(qstrs))
		for _, slice := range qstrs {
			var qslice scp.NodeIDSet
			for _, id := range slice {
				qslice = qslice.Add(scp.NodeID(id))
			}
			q = append(q, qslice)
		}
		node := scp.NewNode(scp.NodeID(nodeID), q, ch, nil)
		nodes[node.ID] = node
		go node.Run(context.Background())
	}

	for slotID := scp.SlotID(1); ; slotID++ {
		msgs := make(map[scp.NodeID]*scp.Msg) // holds the latest message seen from each node

		for _, node := range nodes {
			msgs[node.ID] = nil

			// New slot! Nominate something.
			val := foods[rand.Intn(len(foods))]
			nomMsg := scp.NewMsg(node.ID, slotID, node.Q, &scp.Topic{NomTopic: &scp.NomTopic{X: scp.ValueSet{val}}})
			node.Handle(nomMsg)
		}

		toSend := make(map[scp.NodeID]*scp.Msg)
		for looping := true; looping; {
			select {
			case msg := <-ch:
				if msg.I < slotID {
					// discard messages about old slots
					continue
				}
				msgs[msg.V] = msg
				allExt := true
				for _, m := range msgs {
					if m == nil {
						allExt = false
						break
					}
					if m.T.ExtTopic == nil {
						allExt = false
						break
					}
				}
				if allExt {
					log.Print("all externalized")
					looping = false
					break
				}
				toSend[msg.V] = msg

			default:
				if len(toSend) > 0 {
					for nodeID, msg := range toSend {
						for otherNodeID, otherNode := range nodes {
							if otherNodeID == nodeID {
								continue
							}
							if *delay > 0 {
								otherNode.Delay(rand.Intn(*delay))
							}
							otherNode.Handle(msg)
						}
					}
					toSend = make(map[scp.NodeID]*scp.Msg)
				}
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
