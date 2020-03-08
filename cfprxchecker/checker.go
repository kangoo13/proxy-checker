package cfprxchecker

import (
	"context"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Args struct {
	WebsiteToCrawl        string
	ProxyList             []string
	GoodProxiesOutputFile *os.File
	BadProxiesOutputFile  *os.File
	TimeoutProxy          time.Duration
	muBad                 sync.Mutex
	muGood                sync.Mutex
}

// ProxyURL returns a proxy function (for use in a Transport)
// that always returns the same URL.
func ProxyURL(fixedURL *url.URL) func(*http.Request) (*url.URL, error) {
	return func(pr *http.Request) (*url.URL, error) {
		ctx := context.WithValue(pr.Context(), colly.ProxyURLKey, fixedURL.String())
		*pr = *pr.WithContext(ctx)
		return fixedURL, nil
	}
}

func (a *Args) writeGoodProxy(proxy string) {
	if proxy != "" {

		var b strings.Builder
		a.muGood.Lock()
		defer a.muGood.Unlock()

		b.WriteString(proxy)
		b.WriteString("\n")
		_, err := a.GoodProxiesOutputFile.WriteString(b.String())
		if err != nil {
			log.Fatalf("[WriteGoodProxy] error while appending to file %s", err)
		}
	}
}

func (a *Args) writeBadProxy(proxy string) {
	if proxy != "" {
		var b strings.Builder
		a.muBad.Lock()
		defer a.muBad.Unlock()

		b.WriteString(proxy)
		b.WriteString("\n")
		_, err := a.BadProxiesOutputFile.WriteString(b.String())
		if err != nil {
			log.Fatalf("[WriteBadProxy] error while appending to file %s", err)
		}
	}
}

func NewCollector(args *Args) *colly.Collector {
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

func CheckProxiesAgainstCloudFlare(args *Args) {
	var onErrorCallback = func(response *colly.Response, err error) {
		log.Printf("[DEBUG] Bad proxy found error status %d [%s] [%s]", response.StatusCode, response.Request.ProxyURL, err)
		args.writeBadProxy(response.Request.ProxyURL)
	}

	var onResponseCallback = func(response *colly.Response) {
		log.Printf("[DEBUG] Good proxy found [%s]", response.Request.ProxyURL)
		args.writeGoodProxy(response.Request.ProxyURL)
	}

	var collectors []*colly.Collector

	for _, proxy := range args.ProxyList {
		log.Printf("[DEBUG] Doing %s", proxy)
		newCollector := NewCollector(args)
		newCollector.OnResponse(onResponseCallback)
		newCollector.OnError(onErrorCallback)
		u, err := url.Parse(proxy)
		if err != nil {
			log.Printf("error while parseing proxy %s %s", proxy, err)
			newCollector = nil
			continue
		}
		newCollector.SetProxyFunc(ProxyURL(u))
		if err = newCollector.Visit(args.WebsiteToCrawl); err != nil {
			log.Printf("error happening doing Visit %s", err)
			newCollector = nil
			continue

		}
		collectors = append(collectors, newCollector)
		log.Printf("[DEBUG] end Doing %v", len(collectors))
	}

	for _, collec := range collectors {
		collec.Wait()
	}

}
