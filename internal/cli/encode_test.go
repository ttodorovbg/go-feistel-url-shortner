package cli

import (
	"testing"

	"github.com/ttodorovbg/go-feistel-url-shortener/pkg/codec"
)

func TestRunEncode_FlagKey(t *testing.T) {

	args := []string{
		"--counter", "123",
		"--length", "8",
		"--key", "qweqweqwe",
		"--rounds", "6",
		"--alphabet", codec.Base62Alphabet,
	}

	result, err := runEncode(args)
	if err != nil {
		t.Fatal(err)
	}

	if result != "GoLJBO5W" {
		t.Fatal("expected GoLJBO5W")
	}
}

func TestRunEncode_EnvKey(t *testing.T) {

	t.Setenv("FEISTEL_URL_SHORTENER_KEY", "qweqweqwe")

	args := []string{
		"--counter", "123",
		"--length", "8",
		"--rounds", "6",
		"--alphabet", codec.Base62Alphabet,
	}

	result, err := runEncode(args)
	if err != nil {
		t.Fatal(err)
	}

	if result != "GoLJBO5W" {
		t.Fatal("expected GoLJBO5W")
	}
}

func TestRunEncode_MissingKey(t *testing.T) {

	args := []string{
		"--counter", "1",
	}

	_, err := runEncode(args)

	if err == nil {
		t.Fatal("expected error")
	}
}
