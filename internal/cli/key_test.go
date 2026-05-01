package cli

import (
	"os"
	"testing"
)

func TestGetKey_FlagPriority(t *testing.T) {

	os.Setenv("FEISTEL_URL_SHORTENER_KEY", "env")
	defer os.Unsetenv("FEISTEL_URL_SHORTENER_KEY")

	key, err := getKey("flag")
	if err != nil {
		t.Fatal(err)
	}

	if key != "flag" {
		t.Fatalf("expected flag key, got %s", key)
	}
}

func TestGetKey_FromEnv(t *testing.T) {

	os.Setenv("FEISTEL_URL_SHORTENER_KEY", "env")
	defer os.Unsetenv("FEISTEL_URL_SHORTENER_KEY")

	key, err := getKey("")
	if err != nil {
		t.Fatal(err)
	}

	if key != "env" {
		t.Fatalf("expected env key, got %s", key)
	}
}

func TestGetKey_Error(t *testing.T) {

	os.Unsetenv("FEISTEL_URL_SHORTENER_KEY")

	_, err := getKey("")
	if err == nil {
		t.Fatal("expected error")
	}
}
