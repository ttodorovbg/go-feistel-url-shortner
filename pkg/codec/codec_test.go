package codec_test

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"strconv"
	"sync"
	"testing"

	"github.com/google/uuid"
	c "github.com/ttodorovbg/go-feistel-url-shortner/pkg/codec"
)

func TestCodec_GenerateHash(t *testing.T) {

	type args struct {
		counter uint64
		length  uint8
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"0", args{0, 1}, "D", false},
		{"1", args{1, 1}, "F", false},
		{"2", args{2, 1}, "o", false},
		{"3", args{3071, 2}, "DV", false},
		{"4", args{3072, 2}, "dt", false},
		{"5", args{3073, 2}, "Dh", false},
		{"6", args{1084, 4}, "gMaP", false},
		{"7", args{1085, 4}, "B1XD", false},
		{"8", args{1086, 4}, "9rup", false},
		{"9", args{1087, 5}, "daCJf", false},
		{"10", args{1088, 5}, "f6mxm", false},
		{"11", args{1089, 5}, "54HhA", false},
		{"12", args{1090, 6}, "td4XvJ", false},
		{"13", args{1091, 6}, "l25num", false},
		{"14", args{1092, 6}, "LjkhSC", false},
		{"15", args{1093, 7}, "LdulFT4", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_c := c.NewCodec("test123456789")
			got, err := _c.GenerateHash(tt.args.counter, tt.args.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCodec_ReverseHash(t *testing.T) {
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		args    args
		want    *big.Int
		wantErr bool
	}{
		{"0", args{"D"}, big.NewInt(0), false},
		{"1", args{"F"}, big.NewInt(1), false},
		{"2", args{"o"}, big.NewInt(2), false},
		{"3", args{"DV"}, big.NewInt(3071), false},
		{"4", args{"dt"}, big.NewInt(3072), false},
		{"5", args{"Dh"}, big.NewInt(3073), false},
		{"6", args{"gMaP"}, big.NewInt(1084), false},
		{"7", args{"B1XD"}, big.NewInt(1085), false},
		{"8", args{"9rup"}, big.NewInt(1086), false},
		{"9", args{"daCJf"}, big.NewInt(1087), false},
		{"10", args{"f6mxm"}, big.NewInt(1088), false},
		{"11", args{"54HhA"}, big.NewInt(1089), false},
		{"12", args{"td4XvJ"}, big.NewInt(1090), false},
		{"13", args{"l25num"}, big.NewInt(1091), false},
		{"14", args{"LjkhSC"}, big.NewInt(1092), false},
		{"15", args{"LdulFT4"}, big.NewInt(1093), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_c := c.NewCodec("test123456789")
			got, err := _c.ReverseHash(tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReverseHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReverseHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCodec_GenerateHash_ValidationErrors(t *testing.T) {

	type args struct {
		counter uint64
		length  uint8
		key     []string
		rounds  []uint8
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		err     error
	}{
		{"0", args{1088, 4, []string{}, []uint8{}}, "CWbL", false, nil},
		// rounds
		{"1", args{1088, 4, []string{}, []uint8{0}}, "", true,
			fmt.Errorf("invalid rounds: %d, must be between %d and %d", 0, c.MinRounds(), c.MaxRounds())},
		{"2", args{1088, 4, []string{}, []uint8{2}}, "",
			true, fmt.Errorf("invalid rounds: %d, must be between %d and %d", 2, c.MinRounds(), c.MaxRounds())},
		{"3", args{1088, 4, []string{}, []uint8{3}}, "k8EG", false, nil},
		{"4", args{1088, 4, []string{}, []uint8{10}}, "HnyD", false, nil},
		{"5", args{1088, 4, []string{}, []uint8{11}}, "", true,
			fmt.Errorf("invalid rounds: %d, must be between %d and %d", 11, c.MinRounds(), c.MaxRounds())},
		// length
		{"6", args{1088, 0, []string{}, []uint8{}}, "", true,
			fmt.Errorf("invalid length: %d, must be between %d and %d", 0, 1, c.MaxHashLength())},
		{"7", args{55, 1, []string{}, []uint8{}}, "7", false, nil},
		{"8", args{1088, 12, []string{}, []uint8{}}, "1Xws6JAs85q2", false, nil},
		{"9", args{1088, 11, []string{}, []uint8{}}, "BzblHvJwBwM", false, nil},
		{"10", args{1088, 10, []string{}, []uint8{}}, "XMsA24WCQc", false, nil},
		{"11", args{1088, 9, []string{}, []uint8{}}, "1CdQuSDlO", false, nil},
		{"12", args{1088, 13, []string{}, []uint8{}}, "", true,
			fmt.Errorf("invalid length: %d, must be between %d and %d", 13, 1, c.MaxHashLength())},
		// key
		{"13", args{1088, 10, []string{""}, []uint8{}}, "", true,
			fmt.Errorf("invalid key length: %d, must be between %d and %d", 0, c.MinKeyLength(), c.MaxKeyLength())},
		{"14", args{1088, 10, []string{"Test123"}, []uint8{}}, "", true,
			fmt.Errorf("invalid key length: %d, must be between %d and %d", 7, c.MinKeyLength(), c.MaxKeyLength())},
		{"15", args{1088, 10, []string{"Test1234"}, []uint8{}}, "oBdKMC7sRt", false, nil},
		{"16", args{1088, 10, []string{"SUPER SECRET TEST KEY 01234567891234"}, []uint8{}}, "7PMPi0WK4c", false, nil},
		{"17", args{1088, 10, []string{"SUPER SECRET TEST KEY 012345678912345"}, []uint8{}}, "", true,
			fmt.Errorf("invalid key length: %d, must be between %d and %d", 37, c.MinKeyLength(), c.MaxKeyLength())},
		// counter
		{"18", args{61, 1, []string{}, []uint8{}}, "a", false, nil},
		{"19", args{62, 1, []string{}, []uint8{}}, "", true,
			fmt.Errorf("invalid counter: %d, must be less than or equal %s", 62, strconv.FormatFloat(math.Pow(62, 1)-1, 'f', 0, 64))},
		{"20", args{3843, 2, []string{}, []uint8{}}, "c3", false, nil},
		{"21", args{3844, 2, []string{}, []uint8{}}, "", true,
			fmt.Errorf("invalid counter: %d, must be less than or equal %s", 3844, strconv.FormatFloat(math.Pow(62, 2)-1, 'f', 0, 64))},
		{"22", args{238327, 3, []string{}, []uint8{}}, "qZP", false, nil},
		{"23", args{238328, 3, []string{}, []uint8{}}, "", true,
			fmt.Errorf("invalid counter: %d, must be less than or equal %s", 238328, strconv.FormatFloat(math.Pow(62, 3)-1, 'f', 0, 64))},
		{"24", args{14776335, 4, []string{}, []uint8{}}, "9ZTX", false, nil},
		{"25", args{14776336, 4, []string{}, []uint8{}}, "", true,
			fmt.Errorf("invalid counter: %d, must be less than or equal %s", 14776336, strconv.FormatFloat(math.Pow(62, 4)-1, 'f', 0, 64))},
		{"26", args{916132831, 5, []string{}, []uint8{}}, "1JD5T", false, nil},
		{"27", args{916132832, 5, []string{}, []uint8{}}, "", true,
			fmt.Errorf("invalid counter: %d, must be less than or equal %s", 916132832, strconv.FormatFloat(math.Pow(62, 5)-1, 'f', 0, 64))},
		{"28", args{56800235583, 6, []string{}, []uint8{}}, "Uawe3Z", false, nil},
		{"29", args{56800235584, 6, []string{}, []uint8{}}, "", true,
			fmt.Errorf("invalid counter: %d, must be less than or equal %s", 56800235584, strconv.FormatFloat(math.Pow(62, 6)-1, 'f', 0, 64))},
		{"30", args{3521614606207, 7, []string{}, []uint8{}}, "qh9S9tP", false, nil},
		{"31", args{3521614606208, 7, []string{}, []uint8{}}, "", true,
			fmt.Errorf("invalid counter: %d, must be less than or equal %s", 3521614606208, strconv.FormatFloat(math.Pow(62, 7)-1, 'f', 0, 64))},
		{"32", args{218340105584895, 8, []string{}, []uint8{}}, "sZXQyZ5x", false, nil},
		{"33", args{218340105584896, 8, []string{}, []uint8{}}, "", true,
			fmt.Errorf("invalid counter: %d, must be less than or equal %s", 218340105584896, strconv.FormatFloat(math.Pow(62, 8)-1, 'f', 0, 64))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "test123456789"
			if len(tt.args.key) > 0 {
				key = tt.args.key[0]
			}
			_c := c.NewCodec(key)
			got, err := _c.GenerateHash(tt.args.counter, tt.args.length, tt.args.rounds...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("GenerateHash() error = %v, wantErrMessage %v", err, tt.err)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCodec_ReverseHash_ValidationErrors(t *testing.T) {
	type args struct {
		hash   string
		key    []string
		rounds []uint8
	}
	tests := []struct {
		name    string
		args    args
		want    *big.Int
		wantErr bool
		err     error
	}{
		{"0", args{"abcd", []string{}, []uint8{}}, big.NewInt(3396716), false, nil},
		{"1", args{"abcd&", []string{}, []uint8{}}, nil, true,
			fmt.Errorf("invalid hash character: %c", '&')},
		{"2", args{"abc", []string{""}, []uint8{}}, nil, true,
			fmt.Errorf("invalid key length: %d, must be between %d and %d", len(""), c.MinKeyLength(), c.MaxKeyLength())},
		{"3", args{"abc", []string{"1234567"}, []uint8{}}, nil, true,
			fmt.Errorf("invalid key length: %d, must be between %d and %d", len("1234567"), c.MinKeyLength(), c.MaxKeyLength())},
		{"4", args{"abc", []string{"12345678"}, []uint8{}}, big.NewInt(111967), false, nil},
		{"5", args{"abc", []string{"SUPER SECRET TEST KEY 01234567891234"}, []uint8{}}, big.NewInt(214668), false, nil},
		{"6", args{"abc", []string{"SUPER SECRET TEST KEY 012345678912345"}, []uint8{}}, nil, true,
			fmt.Errorf("invalid key length: %d, must be between %d and %d", 37, c.MinKeyLength(), c.MaxKeyLength())},
		{"7", args{"abc", []string{}, []uint8{0}}, nil, true,
			fmt.Errorf("invalid rounds: %d, must be between %d and %d", 0, c.MinRounds(), c.MaxRounds())},
		{"8", args{"abc", []string{}, []uint8{2}}, nil, true,
			fmt.Errorf("invalid rounds: %d, must be between %d and %d", 2, c.MinRounds(), c.MaxRounds())},
		{"9", args{"abc", []string{}, []uint8{3}}, big.NewInt(147118), false, nil},
		{"10", args{"abc", []string{}, []uint8{10}}, big.NewInt(68860), false, nil},
		{"11", args{"abc", []string{}, []uint8{11}}, nil, true,
			fmt.Errorf("invalid rounds: %d, must be between %d and %d", 11, c.MinRounds(), c.MaxRounds())},
		{"12", args{"", []string{}, []uint8{}}, nil, true,
			fmt.Errorf("invalid length: %d, must be between 1 and %d", 0, c.MaxHashLength())},
		{"13", args{"F", []string{}, []uint8{}}, big.NewInt(1), false, nil},
		{
			"14",
			args{"abcdefghijkl", []string{}, []uint8{}},
			func() *big.Int {
				val, ok := new(big.Int).SetString("1930453752579472361652", 10)
				if !ok {
					panic("failed to parse big.Int string in test")
				}
				return val
			}(),
			false,
			nil,
		},
		{"15", args{"abcdefghijklm", []string{}, []uint8{}}, nil, true, fmt.Errorf("invalid length: %d, must be between 1 and %d", 13, c.MaxHashLength())},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "test123456789"
			if len(tt.args.key) > 0 {
				key = tt.args.key[0]
			}
			_c := c.NewCodec(key)
			got, err := _c.ReverseHash(tt.args.hash, tt.args.rounds...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReverseHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("ReverseHash() error = %v, wantErrMessage %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReverseHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testParams struct {
	maxCount *big.Int
	counts   int
}

func TestCodec_Bidirectional(t *testing.T) {

	tests := make(map[int]testParams)
	tests[1] = testParams{big.NewInt(62 ^ 1 - 1), 61}
	tests[2] = testParams{big.NewInt(62 ^ 21 - 1), 100}
	tests[3] = testParams{big.NewInt(62 ^ 3 - 1), 100}
	tests[4] = testParams{big.NewInt(62 ^ 4 - 1), 100}
	tests[5] = testParams{big.NewInt(62 ^ 5 - 1), 100}
	tests[6] = testParams{big.NewInt(62 ^ 6 - 1), 100}
	tests[7] = testParams{big.NewInt(62 ^ 7 - 1), 100}
	tests[8] = testParams{big.NewInt(62 ^ 8 - 1), 100}
	tests[9] = testParams{big.NewInt(62 ^ 9 - 1), 100}
	tests[10] = testParams{big.NewInt(62 ^ 10 - 1), 100}
	tests[11] = testParams{big.NewInt(62 ^ 11 - 1), 100}
	tests[12] = testParams{big.NewInt(62 ^ 12 - 1), 100}

	uuid := uuid.New().String()

	t.Logf("UUID: %v", uuid)

	var wg sync.WaitGroup

	for u := 8; u <= len(uuid); u++ {
		key := uuid[:u]

		wg.Add(1)
		go func() {
			defer wg.Done()
			testBidirectional(t, key, tests)
		}()
	}

	wg.Wait()
}

func testBidirectional(t *testing.T, key string, tests map[int]testParams) {
	t.Logf("Length of key: %d, Key: %v ", len(key), key)

	_c := c.NewCodec(key)

	for k, v := range tests {
		t.Logf("Test: %d", k)
		// maxHashLength = 12
		for len := uint8(k); len <= 12; len++ {
			for r := c.MinRounds(); r <= c.MaxRounds(); r++ {
				for range v.counts {
					ind := uint64(rand.Int63n(v.maxCount.Int64()))

					hash, err := _c.GenerateHash(ind, len, r)
					if err != nil {
						t.Errorf("Hash() error = %v, key: %v, length: %v, round: %v", err, key, len, r)
						return
					}
					reversed, err := _c.ReverseHash(hash, r)
					if err != nil {
						t.Errorf("ReverseHash() error = %v, key: %v, length: %v, round: %v", err, key, len, r)
						return
					}
					if reversed.Cmp(big.NewInt(int64(ind))) != 0 {
						t.Errorf("Mismatch: counter %d → hash %s → reversed %s, , key: %v, length: %v, round: %v\n",
							ind, hash, reversed.String(), key, len, r)
						break
					}
				}
			}
		}
	}
}

func BenchmarkCodec_GenerateHash(b *testing.B) {
	c := c.NewCodec("super_secret_key")
	for i := 0; b.Loop(); i++ {
		c.GenerateHash(uint64(i), 6)
	}
}
