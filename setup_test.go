package coredns_blackhole

import (
	"coredns/test"
	"testing"

	"github.com/mholt/caddy"
)

func TestSetup(t *testing.T) {
	blocklistFile, rm, err := test.TempFile(".", "google.de")

	if err != nil {
		t.Fatal(err)
	}

	defer rm()

	c := caddy.NewTestController("dns", `blackhole `+blocklistFile)
	if err := setup(c); err != nil {
		t.Fatalf("Unexpected errors, got: %v", err)
	}

	c = caddy.NewTestController("dns", `blackhole`)
	if err := setup(c); err == nil {
		t.Fatalf("Expected errors, but got: %v", err)
	}

	c = caddy.NewTestController("dns", `blackhole invalid`)
	if err := setup(c); err == nil {
		t.Fatalf("Expected errors, but got: %v", err)
	}
}
