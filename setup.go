package coredns_blackhole

import (
	"bufio"
	"fmt"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/mholt/caddy"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	caddy.RegisterPlugin("blackhole", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	blocklist, options, urls, err := parseOptions(c)

	if err != nil {
		return plugin.Error("blocklist", err)
	}

	parseChan := make(chan bool)

	c.OnStartup(func() error {
		go func() {
			ticker := time.NewTicker(options.refresh * time.Second)
			for {
				select {
				case <-parseChan:
					return
				case <-ticker.C:
					i := 0
					for _, url := range urls {
						err := getList(url, blocklist)
						if err != nil {
							log.Debug(err)
						}
						i++
					}
					log.Debugf("Parsed %d lists", i)
				}
			}
		}()
		return nil
	})

	c.OnShutdown(func() error {
		close(parseChan)
		return nil
	})

	dnsserver.GetConfig(c).AddPlugin(
		func(next plugin.Handler) plugin.Handler {
			return Blackhole{
				Next:      next,
				Blocklist: blocklist,
			}
		})

	return nil
}

func parseOptions(c *caddy.Controller) (*Blocklist, *Options, []string, error) {
	blocklist := NewBlocklist()
	options := NewOptions()
	var urls []string
	c.Next() // Skip plugin name

	for c.Next() {
		val := c.Val()
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
			err := getList(val, blocklist)
			if err != nil {
				return nil, nil, nil, err
			}
			urls = append(urls, val)
			log.Debugf("Loaded and parsed blocklist from %s", val)
		}

		finfo, err := os.Stat(val)
		if err == nil && !finfo.IsDir() {
			ret, err := getListsFromFile(val, blocklist, urls)
			if err != nil {
				return nil, nil, nil, err
			}
			urls = append(urls, ret...)
		}

		for c.NextBlock() {
			switch c.Val() {
			case "refresh":
				if !c.NextArg() {
					return nil, nil, nil, c.ArgErr()
				}
				n, err := strconv.Atoi(c.Val())
				if err != nil {
					return nil, nil, nil, err
				}
				if n < 0 {
					return nil, nil, nil, fmt.Errorf("refresh can't be negative: %d", n)
				}
				options.refresh = time.Duration(n)
			case "retry":
				if !c.NextArg() {
					return nil, nil, nil, c.ArgErr()
				}
				n, err := strconv.Atoi(c.Val())
				if err != nil {
					return nil, nil, nil, err
				}
				if n < 0 {
					return nil, nil, nil, fmt.Errorf("retry can't be negative: %d", n)
				}
				options.retry = n
			default:
				return nil, nil, nil, c.Errf("unknown property '%s'", c.Val())
			}

		}
	}

	return blocklist, &options, urls, nil
}

func getListsFromFile(path string, blocklist *Blocklist, urls []string) ([]string, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		url := scanner.Text()
		err := getList(url, blocklist)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func getList(url string, blocklist *Blocklist) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return parseList(resp.Body, blocklist)
}

func parseList(r io.Reader, blocklist *Blocklist) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		var domain string

		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		parts := strings.Fields(line)

		// distinguish between hostfiles and lists of hostnames
		if len(parts) >= 2 {
			if net.ParseIP(parts[0]) != nil {
				match, err := regexp.MatchString(`^(?i:[a-z](?:[a-z0-9\-]*[a-z])?(?:.[a-z](?:[a-z0-9\-]*[a-z])?)*)$`, parts[1])
				if !match || err != nil {
					continue
				}
				domain = parts[1]
			} else {
				continue
			}
		} else {
			domain = parts[0]
		}

		// block local domains
		if strings.HasPrefix("localhost", domain) ||
			domain == "local" ||
			domain == "broadcasthost" ||
			domain == "ip6-localhost" ||
			domain == "ip6-loopback" ||
			domain == "ip6-localnet" ||
			domain == "ip6-mcastprefix" ||
			domain == "ip6-allnodes" ||
			domain == "ip6-allrouters" ||
			domain == "ip6-allhosts" {
			continue

		}

		if !strings.HasSuffix(".", domain) {
			domain = domain + "."
		}

		blocklist.Add(domain)
	}

	return nil
}
