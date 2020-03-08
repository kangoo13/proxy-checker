# proxy-checker

Check for a set of proxies different conditions, is the proxy working, does the proxy bypass cloudflare and so on.
This project is using GoColly as main engine.

# Features

- Checks if the proxy is working
- Checks if the proxy bypass a website protected by CloudFlare

# Installation

Very simple :

```
git clone git@github.com:kangoo13/proxy-checker.git
cd proxy-checker
go build
./proxy-checker -h
Usage: proxy-checker [--timeoutproxy TIMEOUTPROXY] [--cloudflarebypass] [--basic] [--websitecloudflare WEBSITECLOUDFLARE] PROXYLIST GOODPROXIESPATH

Positional arguments:
  PROXYLIST              path to the proxyList
  GOODPROXIESPATH        path to the good proxies identified

Options:
  --timeoutproxy TIMEOUTPROXY, -t TIMEOUTPROXY
                         timeout proxy duration in seconds [default: 4]
  --cloudflarebypass, -c
                         if activated, check if proxy bypass cloudflare
  --basic, -b            if activated, simply checks if proxy is working [default: true]
  --websitecloudflare WEBSITECLOUDFLARE, -w WEBSITECLOUDFLARE
                         website to test proxy against cloudflare, needed if -cf
  --help, -h             display this help and exit

```

# Contribution

Implementing new checkers or enhancing them is welcome !

# Chat

Feel free to speak with me on FreeNode IRC, my nickname is kangoo13
