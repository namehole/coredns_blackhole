package coredns_blackhole

import (
	"bufio"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/mholt/caddy"
	"os"
	"strings"
)

func init() {
	caddy.RegisterPlugin("blackhole", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	c.Next() // 'blackhole'

	files := []string{}
	for c.Next() {
		files = append(files, c.Val())
	}

	if len(files) == 0 {
		return plugin.Error("blackhole", c.ArgErr())
	}

	blocklist := map[string]struct{}{}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return plugin.Error("blackhole", err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix("#", line) {
				if !strings.HasSuffix(".", line) {
					line = line + "."
				}
				blocklist[line] = struct{}{}
			}
		}
	}

	dnsserver.GetConfig(c).AddPlugin(
		func(next plugin.Handler) plugin.Handler {
			return Blackhole{
				Next:      next,
				Blocklist: blocklist,
			}
		})

	return nil
}
