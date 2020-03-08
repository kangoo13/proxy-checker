package checkers

import (
	"context"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"net/http"
	"net/url"
)

// ProxyURL returns a proxy function (for use in a Transport)
// that always returns the same URL.
func ProxyURL(fixedURL *url.URL) func(*http.Request) (*url.URL, error) {
	return func(pr *http.Request) (*url.URL, error) {
		ctx := context.WithValue(pr.Context(), colly.ProxyURLKey, fixedURL.String())
		*pr = *pr.WithContext(ctx)
		return fixedURL, nil
	}
}

func NewCollector(args *Wrapper) *colly.Collector {
	// Instantiate default collector
	c := colly.NewCollector(
		colly.Async(true),
	)

	c.IgnoreRobotsTxt = true
	c.AllowURLRevisit = true
	c.CacheDir = ""

	c.WithTransport(&http.Transport{
		DisableKeepAlives:   true,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	})

	c.SetRequestTimeout(args.TimeoutProxy)

	extensions.RandomUserAgent(c)

	return c
}
