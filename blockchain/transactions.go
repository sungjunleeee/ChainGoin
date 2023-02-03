package blockchain

import (
	"errors"
	"time"

	"github.com/sungjunleeee/ChainGoin/utils"
	"github.com/sungjunleeee/ChainGoin/wallet"
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

// TxIn tracks transaction inputs
type TxIn struct {
	TxID      string `json:"txId"`
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}

// TxOut tracks transaction outputs
type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

// UTxOut tracks unspent transaction outputs
type UTxOut struct {
	TxID   string
	Index  int
	Amount int
}

func (t *Tx) getID() {
	t.ID = utils.Hash(t)
}

func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(t.ID, wallet.Wallet())
	}
}

func validate(tx *Tx) bool {
	valid := true
	for _, txIn := range tx.TxIns {
		// 1. Find the transaction
		// whose unspent outputis used as an input
		prevTx := FindTx(Blockchain(), txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}
		// 2. Get the address (public key) of the unspent output
		address := prevTx.TxOuts[txIn.Index].Address
		// 3. Verify the signature
		valid = wallet.Verify(txIn.Signature, tx.ID, address)
		if !valid {
			break
		}
	}
	return valid
}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
Outer:
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer
			}
		}
	}
	return exists
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

var (
	ErrorNotEnoughBalance = errors.New("Not enough balance")
	ErrorNotValid         = errors.New("Not valid transaction")
)

func makeTx(from, to string, amount int) (*Tx, error) {
	if GetBalanceByAddress(from, Blockchain()) < amount {
		return nil, ErrorNotEnoughBalance
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := FilterUTxOutsByAddress(from, Blockchain())
	// 1. Find enough unspent TxOuts for transfer
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{
			TxID:      uTxOut.TxID,
			Index:     uTxOut.Index,
			Signature: from,
		}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	// 2. Create TxOuts if change is needed
	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{
			Address: from,
			Amount:  change,
		}
		txOuts = append(txOuts, changeTxOut)
	}
	// 3. Create TxOuts for transfer
	txOut := &TxOut{
		Address: to,
		Amount:  amount,
	}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getID()
	tx.sign()
	valid := validate(tx)
	if !valid {
		return nil, ErrorNotValid
	}
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return ErrorNotValid
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) LoadMempool() []*Tx {
	coinbase := makeCoinbaseTx(wallet.Wallet().Address) // reward to the miner
	txs := append(m.Txs, coinbase)
	m.Txs = nil // reset mempool
	return txs
}
