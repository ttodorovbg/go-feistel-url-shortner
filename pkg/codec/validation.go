package codec

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

const MinRounds uint8 = 3
const MaxRounds uint8 = 10
const MinHashLength uint8 = 1
const MaxHashLength uint8 = 12
const MinKeyLength uint8 = 8
const MaxKeyLength uint8 = 36 // can use UUID

type ErrInvalidRounds struct {
	Value uint8
	Min   uint8
	Max   uint8
}

type ErrInvalidHashChar struct {
	Char rune
}

type ErrInvalidHashLength struct {
	Length uint8
	Min    uint8
	Max    uint8
}

type ErrInvalidKeyLength struct {
	Length uint8
	Min    uint8
	Max    uint8
}

type ErrInvalidCounter struct {
	Counter big.Int
	Min     big.Int
	Max     big.Int
}

func (e *ErrInvalidRounds) Error() string {
	return fmt.Sprintf("invalid rounds: %d, must be between %d and %d", e.Value, e.Min, e.Max)
}

func (e *ErrInvalidHashChar) Error() string {
	return fmt.Sprintf("invalid hash character: %c", e.Char)
}

func (e *ErrInvalidHashLength) Error() string {
	return fmt.Sprintf("invalid hash length: %d, must be between %d and %d", e.Length, e.Min, e.Max)
}

func (e *ErrInvalidKeyLength) Error() string {
	return fmt.Sprintf("invalid key length: %d, must be between %d and %d", e.Length, e.Min, e.Max)
}

func (e *ErrInvalidCounter) Error() string {
	return fmt.Sprintf("invalid counter: %s, must be between %s and %s", e.Counter.String(), e.Min.String(), e.Max.String())
}

func validateRounds(rounds uint8) error {
	if rounds < MinRounds || rounds > MaxRounds {
		return &ErrInvalidRounds{
			Value: rounds,
			Min:   MinRounds,
			Max:   MaxRounds,
		}
	}
	return nil
}

func validateHashChars(hash string) error {
	for _, c := range hash {
		if !strings.ContainsRune(base62Alphabet, c) {
			return &ErrInvalidHashChar{Char: c}
		}
	}

	return nil
}

func validateHashLength(length uint8) error {
	if length < MinHashLength || length > MaxHashLength {
		return &ErrInvalidHashLength{
			Length: length,
			Min:    MinHashLength,
			Max:    MaxHashLength,
		}
	}
	return nil
}

func validateKeyLength(key string) error {
	if len(key) < int(MinKeyLength) || len(key) > int(MaxKeyLength) {
		return &ErrInvalidKeyLength{
			Length: uint8(len(key)),
			Min:    MinKeyLength,
			Max:    MaxKeyLength,
		}
	}
	return nil
}

func validateCounter(counter uint64, length uint8) error {
	maxCounter := computeMaxBase62(int64(length))

	counterBigInt := new(big.Int).SetUint64(counter)
	if maxCounter.Cmp(big.NewInt(0)) < 0 || counterBigInt.Cmp(maxCounter) > 0 {
		return &ErrInvalidCounter{
			Counter: *counterBigInt,
			Min:     *big.NewInt(0),
			Max:     *maxCounter,
		}
	}

	return nil
}

func computeMaxBase62(length int64) *big.Int {

	if length <= int64(MaxLengthForFloat64) {
		f := math.Pow(float64(Base62), float64(length)) - 1
		return big.NewInt(int64(f))
	}
	exponent := big.NewInt(length)
	max := new(big.Int).Exp(big62, exponent, nil)
	return max.Sub(max, big.NewInt(1))
}
