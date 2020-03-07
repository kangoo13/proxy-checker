package cfprxchecker

import (
	"context"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"net/url"
	"sync/atomic"
)

type roundRobinSwitcher struct {
	proxyURLs    []*url.URL
	proxyURLsMap map[string]*url.URL
	index        uint32
}

func (r *roundRobinSwitcher) GetProxy(pr *http.Request) (*url.URL, error) {
	var u *url.URL
	// Use Same Proxy when Redirect
	if pr.Response != nil && pr.Response.StatusCode == 301 {
		u = r.proxyURLsMap[pr.Response.Request.Context().Value(colly.ProxyURLKey).(string)]
		log.Printf("titi %s", u.String())
	} else if pr.Response != nil {
		log.Printf("tototototo")
	} else {
		u = r.proxyURLs[r.index%uint32(len(r.proxyURLs))]
		atomic.AddUint32(&r.index, 1)
		log.Printf("grosminet %s", u.String())
	}

	ctx := context.WithValue(pr.Context(), colly.ProxyURLKey, u.String())
	*pr = *pr.WithContext(ctx)

	return u, nil
}

// RoundRobinProxySwitcher creates a proxy switcher function which rotates
// ProxyURLs on every request.
// The proxy type is determined by the URL scheme. "http", "https"
// and "socks5" are supported. If the scheme is empty,
// "http" is assumed.
func RoundRobinProxySwitcher(ProxyURLs ...string) (colly.ProxyFunc, error) {
	urls := make([]*url.URL, len(ProxyURLs))
	urlsMap := make(map[string]*url.URL, len(ProxyURLs))
	for i, u := range ProxyURLs {
		parsedU, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		urls[i] = parsedU
		urlsMap[parsedU.String()] = parsedU
	}
	return (&roundRobinSwitcher{
		proxyURLs:    urls,
		proxyURLsMap: urlsMap,
		index:        0,
	}).GetProxy, nil
}
