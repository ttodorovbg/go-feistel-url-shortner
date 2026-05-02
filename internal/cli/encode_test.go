package cli

import (
	"testing"
)

func TestRunEncode_FlagKey(t *testing.T) {

	args := []string{
		"--counter", "123",
		"--length", "8",
		"--key", "qweqweqwe",
		"--rounds", "6",
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
