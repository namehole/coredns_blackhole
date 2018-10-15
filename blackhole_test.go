package coredns_blackhole

import (
	"context"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
	"testing"
)

func TestBlackhole_ServeDNS(t *testing.T) {
	// Create a new Example Plugin. Use the test.ErrorHandler as the next plugin.
	list := NewBlocklist()
	list.Add("example.net.")
	bl := Blackhole{Next: test.ErrorHandler(), Blocklist: list}

	ctx := context.TODO()

	tests := []struct {
		reqDomain     string
		shouldErr     bool
		expectedRCode int
	}{
		{
			"example.net",
			false,
			BlockedRcode,
		},

		{
			"example.com",
			false,
			dns.RcodeServerFailure,
		},
	}

	for i, testcase := range tests {
		req := new(dns.Msg)

		req.SetQuestion(dns.Fqdn(testcase.reqDomain), dns.TypeA)

		rec := dnstest.NewRecorder(&test.ResponseWriter{})

		code, err := bl.ServeDNS(ctx, rec, req)

		if err == nil && testcase.shouldErr {
			t.Fatalf("Test %d: expected errors, but got no error", i)
		} else if err != nil && !testcase.shouldErr {
			t.Fatalf("Test %d: expected no errors, but got '%v'", i, err)
			continue
		}

		if code != testcase.expectedRCode {
			t.Errorf("Test %d: Expected status code %d, but got %d", i, testcase.expectedRCode, code)
		}

		if rec.Rcode != testcase.expectedRCode {
			t.Errorf("Test %d: Expected return code %d, but got %d", i, testcase.expectedRCode, rec.Rcode)
		}
	}
}
