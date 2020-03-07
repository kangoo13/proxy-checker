# cloudflare-proxy-ban-checker

Checks a set of proxies against a protected CloudFlare website to check if the proxies are banned or not.

# Installation

Very simple :

```
git clone git@github.com:kangoo13/cloudflare-proxy-ban-checker.git
cd cloudflare-proxy-ban-checker
go build
./cloudflare-proxy-ban-checker -h
Usage: cloudflare-proxy-ban-checker [--goodproxiespath GOODPROXIESPATH] [--badproxiespath BADPROXIESPATH] [--timeoutproxy TIMEOUTPROXY] WEBSITE PROXYLIST

Positional arguments:
  WEBSITE                website to test proxy against cloudflare
  PROXYLIST              path to the proxyList

Options:
  --goodproxiespath GOODPROXIESPATH
                         path to the good proxies identified [default: good.txt]
  --badproxiespath BADPROXIESPATH
                         path to the bad proxies identified [default: bad.txt]
  --timeoutproxy TIMEOUTPROXY
                         timeout proxy duration [default: 5]
  --help, -h             display this help and exit
```

