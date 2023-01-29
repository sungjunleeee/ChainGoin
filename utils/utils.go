package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

func HandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ToByte(data interface{}) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	HandleErr(encoder.Encode(data))
	return buffer.Bytes()
}
