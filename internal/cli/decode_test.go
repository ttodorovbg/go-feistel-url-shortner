package cli

import (
	"os"
	"testing"

	"github.com/ttodorovbg/go-feistel-url-shortener/pkg/codec"
)

func TestRunDecode_FlagKey(t *testing.T) {

	hash, err := codec.GenerateHash(123, 8, "qweqweqwe", 6)
	if err != nil {
		t.Fatal(err)
	}

	args := []string{
		"--hash", hash,
		"--key", "qweqweqwe",
	}

	result, err := runDecode(args)
	if err != nil {
		t.Fatal(err)
	}

	if result != "123" {
		t.Fatalf("expected 123, got %s", result)
	}
}

func TestRunDecode_EnvKey(t *testing.T) {

	os.Setenv("FEISTEL_URL_SHORTENER_KEY", "qweqweqwe")
	defer os.Unsetenv("FEISTEL_URL_SHORTENER_KEY")

	hash, err := codec.GenerateHash(123, 8, "qweqweqwe", 6)
	if err != nil {
		t.Fatal(err)
	}

	args := []string{
		"--hash", hash,
	}

	result, err := runDecode(args)
	if err != nil {
		t.Fatal(err)
	}

	if result != "123" {
		t.Fatalf("expected 123, got %s", result)
	}
}

func TestRunDecode_MissingHash(t *testing.T) {

	args := []string{
		"--key", "secret",
	}

	_, err := runDecode(args)

	if err == nil {
		t.Fatal("expected error")
	}
}
