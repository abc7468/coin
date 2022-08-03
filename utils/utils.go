package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"strings"
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

func FromBytes(x any, data []byte) {
	encoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(encoder.Decode(x))
}

func Hash(i any) string {
	s := fmt.Sprintf("%v", i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func GetSplitedStrings(s string, sep string, i int) string {
	res := strings.Split(s, sep)
	if i >= len(res) {
		return ""
	}
	return res[i]
}
