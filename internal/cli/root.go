package cli

import (
	"errors"
	"fmt"
	"os"
)

func Execute() (string, error) {

	if len(os.Args) < 2 {
		return "", errors.New("expected 'encode' or 'decode'")
	}

	switch os.Args[1] {
	case "encode":
		return runEncode(os.Args[2:])
	case "decode":
		return runDecode(os.Args[2:])
	default:
		return "", fmt.Errorf("unknown command: %s", os.Args[1])
	}
}
