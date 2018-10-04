package coredns_blackhole

import (
	"context"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
	"testing"
)

func TestExample(t *testing.T) {
	// Create a new Example Plugin. Use the test.ErrorHandler as the next plugin.
	bl := Blackhole{Next: nil, Blocklist: map[string]struct{}{
			"google.com.": struct{}{},
			"facebook.com.": struct {}{},
		},
	}

	ctx := context.TODO()

	tests := []struct {
		reqDomain    string
		shouldErr    bool
		shouldRefuse bool
	}{
		{
			"google.com",
			false,
			true,
		},

		{
			"example.com",
			true,
			false,
		},
	}

	for i, testing := range tests {
		req := new(dns.Msg)

		req.SetQuestion(dns.Fqdn(testing.reqDomain), dns.TypeA)

		rec := dnstest.NewRecorder(&test.ResponseWriter{})

		code, err := bl.ServeDNS(ctx, rec, req)

		if err == nil && testing.shouldErr {
			t.Fatalf("Test %d expected errors, but got no error", i)
		} else if err != nil && !testing.shouldErr {
			t.Fatalf("Test %d expected no errors, but got '%v'", i, err)
		}

		if code != dns.RcodeRefused && testing.shouldRefuse {
			t.Errorf("Expected status code %d, but got %d", dns.RcodeRefused, code)
		}
	}
}
