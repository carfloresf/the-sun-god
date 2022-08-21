package random

import (
	"crypto/rand"
	"math/big"
)

var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandSequence(seed int) string {
	b := make([]rune, seed)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return ""
		}

		b[i] = letters[n.Int64()]
	}

	return string(b)
}
