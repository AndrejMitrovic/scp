package main

import (
	"github.com/bobg/scp"

	"github.com/chain/txvm/protocol/bc"
)

func nominate() {
	txpool := make(map[bc.Hash]*bc.Tx)

	doNom := func() error {
		if len(txpool) == 0 {
			return nil
		}

		txs := make([]*bc.Tx, 0, len(txpool))
		for _, tx := range txpool {
			txs = append(txs, tx)
		}

		block, _, err := chain.GenerateBlock(ctx, chain.State(), timestampMS, txs)
		if err != nil {
			return err
		}
		// xxx figure out which txs GenerateBlock removed as invalid, and remove them from txpool

		err = storeBlock(block)
		if err != nil {
			return err
		}
		n := &scp.NomTopic{
			X: scp.ValueSet{block.Hash()},
		}
		msg := scp.NewMsg(node.ID, scp.SlotID(block.Height), node.Q, n) // xxx slotID is 32 bits, block height is 64
		node.Handle(msg)
	}

	for item := range nomChan {
		switch item := item.(type) {
		case *bc.Tx:
			if _, ok := txpool[item.ID]; ok {
				// tx is already in the pool
				continue
			}
			txpool[item.ID] = item // xxx need to persist this
			err := doNom()
			if err != nil {
				panic(err) // xxx
			}

		case scp.SlotID:
			err := doNom()
			if err != nil {
				panic(err) // xxx
			}
		}
	}
}
