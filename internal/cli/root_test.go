package cli

import (
	"os"
	"testing"

	"github.com/ttodorovbg/go-feistel-url-shortener/pkg/codec"
)

func TestExecute_Encode(t *testing.T) {

	os.Args = []string{
		"cmd",
		"encode",
		"--counter", "123",
		"--length", "8",
		"--key", "qweqweqwe",
	}

	result, err := Execute()
	if err != nil {
		t.Fatal(err)
	}

	if result != "GoLJBO5W" {
		t.Fatal("expected GoLJBO5W")
	}
}

func TestExecute_Decode(t *testing.T) {

	hash, err := codec.GenerateHash(123, 8, "qweqweqwe", 6)
	if err != nil {
		t.Fatal(err)
	}

	os.Args = []string{
		"cmd",
		"decode",
		"--hash", hash,
		"--key", "qweqweqwe",
	}

	result, err := Execute()
	if err != nil {
		t.Fatal(err)
	}

	if result != "123" {
		t.Fatalf("expected 123, got %s", result)
	}
}

func TestExecute_UnknownCommand(t *testing.T) {

	os.Args = []string{"cmd", "unknown"}

	_, err := Execute()

	if err == nil {
		t.Fatal("expected error")
	}
}
