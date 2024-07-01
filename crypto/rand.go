package crypto

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

func init() {
	assertAvailablePRNG()
}

const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
)

func assertAvailablePRNG() {
	// Assert that a cryptographically secure PRNG is available.
	// Panic otherwise.
	buf := make([]byte, 1)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: Read() failed with %#v", err))
	}
}

// RandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func RandomString(n int) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		ret[i] = alphabet[num.Int64()]
	}

	return string(ret), nil
}
