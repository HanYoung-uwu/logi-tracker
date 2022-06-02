package utility

import (
	"crypto/rand"
	"log"
	"math/big"
)

func RandBytes(l int) []byte {
	store := "asdfghjklqwertyuiop1234567890zxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	maxLen := len(store)
	result := make([]byte, l)
	i := 0
	for {
		if i == l {
			break
		}
		p, err := rand.Int(rand.Reader, big.NewInt(int64(maxLen)))
		if err != nil {
			log.Fatal(err)
		}
		result[i] = store[p.Int64()]
		i += 1
	}
	return result
}
