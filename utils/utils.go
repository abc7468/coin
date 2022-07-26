package utils

import (
	"bytes"
	"encoding/gob"
)

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ToBytes(x any) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	HandleErr(encoder.Encode(x))
	return aBuffer.Bytes()
}
