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

func (t *Tx) getId() {
	t.ID = utils.Hash(t)
}

type TxIn struct {
	Owner  string
	Amount int
}

type TxOut struct {
	Owner  string
	Amount int
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"COINBASE", minerReward},
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
	tx.getId()
	return tx
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if Blockchain().GetBalanceByAddress(from) < amount {
		return nil, errors.New("Not enough balance")
	}
	var txIns []*TxIn
	var txOuts []*TxOut
	total := 0
	prevTxOuts := Blockchain().FilterTxOutsByAddress(from)
	for _, txOut := range prevTxOuts {
		if total >= amount {
			break
		}
		txIns = append(txIns, &TxIn{txOut.Owner, txOut.Amount})
		total += txOut.Amount
	}
	change := total - amount
	if change != 0 {
		// Change back to the sender
		txOuts = append(txOuts, &TxOut{from, change})
	}
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
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
