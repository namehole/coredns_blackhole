package coredns_blackhole

import (
	"context"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin("blackhole")

type Blackhole struct {
	Next      plugin.Handler
	Blocklist map[string]struct{}
}

func (b Blackhole) Name() string { return "blackhole" }

func (b Blackhole) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	_, ok := b.Blocklist[state.Name()]
	if ok {
		return dns.RcodeRefused, nil
	}

	return plugin.NextOrFailure(b.Name(), b.Next, ctx, w, r)
}
