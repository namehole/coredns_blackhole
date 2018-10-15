[![Build Status](https://travis-ci.com/namehole/coredns_blackhole.svg?branch=master)](https://travis-ci.com/namehole/coredns_blackhole)

# coredns_blackhole

## Name

_blackhole_ - Applies bocklists to DNS queries

## Description

The _blackhole_ plugin downloads blocklists and applies them to the DNS queries. If a requested address matches a blocklist entry, the request will be rejected.

## Syntax

```
blackhole [BLOCKLIST_FILE... ] [BLOCKLIST_URL... ] {
    refresh SECONDS
}
```

If the argument is a valid path to a file, the file will be interpreted as a list of blocklist urls. Each url is then downloaded and added to the blocklist.

If the argument is a valid url, the url is downloaded and added to the blocklist.

The `refresh` option sets the timer for the refresh of the blocklists. Default is 30 seconds.

## Examples

Block all urls in the simple_ad list and forward the rest.

```
. {
    blackhole https://s3.amazonaws.com/lists.disconnect.me/simple_ad.txt {
        refresh 60
    }
    forward tls://1.1.1.1 tls://1.0.0.1
}
    
```
