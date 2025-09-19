package codec

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"strings"
)

// https://www.ietf.org/rfc/rfc3986.txt
/*

   Characters that are allowed in a URI but do not have a reserved
   purpose are called unreserved.  These include uppercase and lowercase
   letters, decimal digits, hyphen, period, underscore, and tilde.

      unreserved  = ALPHA / DIGIT / "-" / "." / "_" / "~"
*/

type Codec struct {
	key    string
	rounds uint8
}

const base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const defaultRounds uint8 = 6

var big62 = big.NewInt(62)

func NewCodec(key string, args ...uint8) *Codec {
	rounds := defaultRounds
	if len(args) > 0 {
		rounds = args[0]
	}
	return &Codec{key: key, rounds: rounds}
}

func toBase62(num *big.Int, length uint8) string {
	n := new(big.Int).Set(num)
	chars := make([]byte, length)
	mod := new(big.Int)

	for i := int(length) - 1; i >= 0; i-- {
		n.DivMod(n, big62, mod)
		chars[i] = base62Alphabet[mod.Int64()]
	}
	return string(chars)
}

func fromBase62(s string) *big.Int {
	num := big.NewInt(0)
	for _, c := range s {
		index := int64(strings.IndexRune(base62Alphabet, c))
		num.Mul(num, big62)
		num.Add(num, big.NewInt(index))
	}
	return num
}

func roundFunction(right *big.Int, roundKey string, mod *big.Int) *big.Int {
	data := fmt.Appendf(nil, "%s:%s", right.String(), roundKey)
	hash := sha256.Sum256(data)
	val := binary.BigEndian.Uint32(hash[0:4])
	return new(big.Int).Mod(big.NewInt(int64(val)), mod)
}

func feistelBijective(counter *big.Int, rounds uint8, key string, domain *big.Int) *big.Int {
	bitLen := domain.BitLen()
	halfSize := bitLen / 2

	leftMask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), uint(halfSize)), big.NewInt(1))
	rightMask := new(big.Int).Set(leftMask)

	R := new(big.Int).And(counter, rightMask)
	L := new(big.Int).Rsh(counter, uint(halfSize))
	L.And(L, leftMask)

	mod := new(big.Int).Lsh(big.NewInt(1), uint(halfSize))

	for i := range rounds {
		F := roundFunction(R, fmt.Sprintf("%s-%d", key, i), mod)
		tmp := new(big.Int).Xor(L, F)
		L, R = R, tmp
	}

	combined := new(big.Int).Or(new(big.Int).Lsh(L, uint(halfSize)), R)
	if combined.Cmp(domain) >= 0 {
		return feistelBijective(combined, rounds, key, domain)
	}
	return combined
}

func feistelInverse(scrambled *big.Int, rounds uint8, key string, domain *big.Int) *big.Int {
	bitLen := domain.BitLen()
	halfSize := bitLen / 2

	leftMask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), uint(halfSize)), big.NewInt(1))
	rightMask := new(big.Int).Set(leftMask)

	R := new(big.Int).And(scrambled, rightMask)
	L := new(big.Int).Rsh(scrambled, uint(halfSize))
	L.And(L, leftMask)

	mod := new(big.Int).Lsh(big.NewInt(1), uint(halfSize))

	for i := int(rounds) - 1; i >= 0; i-- {
		F := roundFunction(L, fmt.Sprintf("%s-%d", key, i), mod)
		tmp := new(big.Int).Xor(R, F)
		R, L = L, tmp
	}

	combined := new(big.Int).Or(new(big.Int).Lsh(L, uint(halfSize)), R)
	if combined.Cmp(domain) >= 0 {
		return feistelInverse(combined, rounds, key, domain) // cycle-walk
	}
	return combined
}

func GenerateHash(counter uint64, length uint8, key string, args ...uint8) (string, error) {

	rounds := defaultRounds

	if len(args) > 0 {
		rounds = args[0]
	}

	if err := validateRounds(rounds); err != nil {
		return "", err
	}
	if err := validateLength(length); err != nil {
		return "", err
	}
	if err := validateKey(key); err != nil {
		return "", err
	}
	if err := validateCounter(counter, length); err != nil {
		return "", err
	}

	var domain *big.Int

	if length <= 8 {
		f := math.Pow(62, float64(length))
		domain = big.NewInt(int64(f))
	} else {
		domain = new(big.Int).Exp(big62, big.NewInt(int64(length)), nil)
	}

	scrambled := feistelBijective(big.NewInt(int64(counter)), rounds, key, domain)

	return toBase62(scrambled, length), nil
}

func ReverseHash(hash string, key string, args ...uint8) (*big.Int, error) {

	if err := validateHash(hash); err != nil {
		return nil, err
	}
	if err := validateKey(key); err != nil {
		return nil, err
	}

	rounds := defaultRounds

	if len(args) > 0 {
		rounds = args[0]
	}

	length := len(hash)

	if err := validateRounds(rounds); err != nil {
		return nil, err
	}
	if err := validateLength(uint8(length)); err != nil {
		return nil, err
	}

	var domain *big.Int

	if length <= 8 {
		f := math.Pow(62, float64(length))
		domain = big.NewInt(int64(f))
	} else {
		domain = new(big.Int).Exp(big62, big.NewInt(int64(length)), nil)
	}

	scrambled := fromBase62(hash)

	return feistelInverse(scrambled, rounds, key, domain), nil
}

func (c *Codec) GenerateHash(counter uint64, length uint8, args ...uint8) (string, error) {
	return GenerateHash(counter, length, c.key, args...)
}

func (c *Codec) ReverseHash(hash string, args ...uint8) (*big.Int, error) {
	return ReverseHash(hash, c.key, args...)
}
