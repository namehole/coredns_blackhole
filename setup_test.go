package coredns_blackhole

import (
	"github.com/coredns/coredns/test"
	"testing"

	"github.com/mholt/caddy"
)

func TestSetup_BlocklistURLs(t *testing.T) {
	blocklistURLs := "https://s3.amazonaws.com/lists.disconnect.me/simple_ad.txt https://s3.amazonaws.com/lists.disconnect.me/simple_tracking.txt"

	c := caddy.NewTestController("dns", `blackhole `+blocklistURLs)
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

func TestSetup_BlocklistFile(t *testing.T) {
	blocklistFile, rm, err := test.TempFile(".", "https://s3.amazonaws.com/lists.disconnect.me/simple_ad.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer rm()

	c := caddy.NewTestController("dns", `blackhole `+blocklistFile)
	if err := setup(c); err != nil {
		t.Fatalf("Unexpected errors, got: %v", err)
	}
}

func TestSetup_Invalid(t *testing.T) {
	c := caddy.NewTestController("dns", `blackhole`)
	if err := setup(c); err == nil {
		t.Fatalf("Expected errors, but got: %v", err)
	}

	c = caddy.NewTestController("dns", `blackhole invalid`)
	if err := setup(c); err == nil {
		t.Fatalf("Expected errors, but got: %v", err)
	}
}
