package codec

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

const minRounds uint8 = 3
const maxRounds uint8 = 10
const maxHashLength uint8 = 12
const minKeyLength uint8 = 8
const maxKeyLength uint8 = 36 // can use UUID

func validateRounds(rounds uint8) error {
	if rounds < minRounds || rounds > maxRounds {
		return fmt.Errorf("invalid rounds: %d, must be between %d and %d", rounds, minRounds, maxRounds)
	}
	return nil
}

func validateLength(length uint8) error {
	if length < 1 || length > maxHashLength {
		return fmt.Errorf("invalid length: %d, must be between 1 and %d", length, maxHashLength)
	}
	return nil
}

func validateKey(key string) error {
	if len(key) < int(minKeyLength) || len(key) > int(maxKeyLength) {
		return fmt.Errorf("invalid key length: %d, must be between %d and %d", len(key), minKeyLength, maxKeyLength)
	}
	return nil
}

func validateHash(hash string) error {
	for _, c := range hash {
		if !strings.ContainsRune(base62Alphabet, c) {
			return fmt.Errorf("invalid hash character: %c", c)
		}
	}
	return nil
}

func computeMaxBase62(length int64) *big.Int {

	if length <= int64(MaxLengthForFloat64) {
		// Safe conversion from float64 to *big.Int for small exponents
		f := math.Pow(float64(Base62), float64(length)) - 1
		return big.NewInt(int64(f))
	}
	exponent := big.NewInt(length)
	max := new(big.Int).Exp(big62, exponent, nil)
	return max.Sub(max, big.NewInt(1))
}

func validateCounter(counter uint64, length uint8) error {
	maxCounter := computeMaxBase62(int64(length))

	counterBigInt := new(big.Int).SetUint64(counter)
	if counterBigInt.Cmp(maxCounter) > 0 {
		return fmt.Errorf("invalid counter: %d, must be less than or equal %s", counter, maxCounter.String())
	}

	return nil
}

func MinRounds() uint8 {
	return minRounds
}

func MaxRounds() uint8 {
	return maxRounds
}

func MaxHashLength() uint8 {
	return maxHashLength
}

func MinKeyLength() uint8 {
	return minKeyLength
}

func MaxKeyLength() uint8 {
	return maxKeyLength
}
