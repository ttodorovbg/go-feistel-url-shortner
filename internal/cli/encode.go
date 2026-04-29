package cli

import (
	"flag"

	"github.com/ttodorovbg/go-feistel-url-shortner/pkg/codec"
)

func runEncode(args []string) (string, error) {

	// cmd := flag.NewFlagSet("encode", flag.ExitOnError)
	cmd := flag.NewFlagSet("encode", flag.ErrorHandling(flag.ContinueOnError))

	counter := cmd.Uint64("counter", 0, "")
	length := cmd.Uint("length", 8, "")
	keyFlag := cmd.String("key", "", "secret key (optional if FEISTEL_URL_SHORTENER_KEY env is set)")
	rounds := cmd.Uint("rounds", 6, "rounds")

	cmd.Parse(args)

	key, err := getKey(*keyFlag)
	if err != nil {
		return "", err
	}

	hash, err := codec.GenerateHash(uint64(*counter), uint8(*length), key, uint8(*rounds))
	if err != nil {
		return "", err
	}

	return hash, nil
}
