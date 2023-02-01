package blockchain

import (
	"errors"
	"time"

	"github.com/sungjunleeee/ChainGoin/utils"
)

const (
	minerReward = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (t *Tx) getID() {
	t.ID = utils.Hash(t)
}

// TxIn tracks transaction inputs
type TxIn struct {
	TxID  string `json:"txId"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

// TxOut tracks transaction outputs
type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

// UTxOut tracks unspent transaction outputs
type UTxOut struct {
	TxID   string
	Index  int
	Amount int
}

func isOnMempool(uTxOut *UTxOut) bool {
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			return input.TxID == uTxOut.TxID && input.Index == uTxOut.Index
		}
	}
	return false
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getID()
	return tx
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if Blockchain().GetBalanceByAddress(from) < amount {
		return nil, errors.New("Not enough balance")
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := Blockchain().FilterUTxOutsByAddress(from)
	// 1. Find enough unspent TxOuts for transfer
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{
			TxID:  uTxOut.TxID,
			Index: uTxOut.Index,
			Owner: from,
		}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	// 2. Create TxOuts if change is needed
	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{
			Owner:  from,
			Amount: change,
		}
		txOuts = append(txOuts, changeTxOut)
	}
	// 3. Create TxOuts for transfer
	txOut := &TxOut{
		Owner:  to,
		Amount: amount,
	}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getID()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("Jun", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) LoadMempool() []*Tx {
	coinbase := makeCoinbaseTx("Jun") // reward to the miner
	txs := append(m.Txs, coinbase)
	m.Txs = nil // reset mempool
	return txs
}
