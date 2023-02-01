package db

import (
	"github.com/boltdb/bolt"
	"github.com/sungjunleeee/ChainGoin/utils"
)

const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		newDB, err := bolt.Open(dbName, 0600, nil)
		db = newDB
		utils.HandleErr(err)
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = tx.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

func Close() {
	DB().Close()
}

func SaveBlock(hash string, data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func SaveBlockchain(data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dataBucket))
		err := b.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleErr(err)
}

func LoadBlockchain() []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(dataBucket))
		data = b.Get([]byte(checkpoint))
		return nil
	})
	return data
}

func Block(hash string) []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		data = b.Get([]byte(hash))
		return nil
	})
	return data
}
