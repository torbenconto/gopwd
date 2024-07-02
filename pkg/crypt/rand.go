package crypt

import (
	cryptoRandom "crypto/rand"
	"math/big"
	random "math/rand"
)

func RandomInt(max int) int {
	rand, err := cryptoRandom.Int(cryptoRandom.Reader, big.NewInt(int64(max)))
	if err != nil {
		// Fall back to math/rand if crypto/rand fails
		return random.Intn(max)
	}

	return int(rand.Int64())
}
