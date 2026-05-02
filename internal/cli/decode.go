package cli

import (
	"errors"
	"flag"

	"github.com/ttodorovbg/go-feistel-url-shortener/pkg/codec"
)

func runDecode(args []string) (string, error) {

	cmd := flag.NewFlagSet("decode", flag.ExitOnError)

	hash := cmd.String("hash", "", "hash string")
	keyFlag := cmd.String("key", "", "secret key (optional if FEISTEL_URL_SHORTENER_KEY env is set)")
	rounds := cmd.Uint("rounds", 6, "rounds")

	if err := cmd.Parse(args); err != nil {
		return "", err
	}

	if *hash == "" {
		return "", errors.New("--hash is required")
	}

	key, err := getKey(*keyFlag)
	if err != nil {
		return "", err
	}

	counter, err := codec.ReverseHash(*hash, key, uint8(*rounds))
	if err != nil {
		return "", err
	}

	return string(counter.Text(10)), nil
}
