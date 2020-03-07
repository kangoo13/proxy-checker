package cfprxchecker

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"log"
	"net/http"
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

func CheckProxiesAgainstCloudFlare(args *Args) {
	// Rotate the proxies
	rp, err := RoundRobinProxySwitcher(args.ProxyList...)
	if err != nil {
		log.Fatal(err)
	}

	// Instantiate default collector
	c := colly.NewCollector(
		colly.Async(true),
	)

	c.IgnoreRobotsTxt = true
	c.AllowURLRevisit = true
	c.CacheDir = ""

	c.WithTransport(&http.Transport{
		Proxy:               rp,
		DisableKeepAlives:   true,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	})

	// Limit the maximum parallelism to 24
	err = c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 24})

	if err != nil {
		log.Printf("error while doing c.Limit %s", err)
	}

	extensions.RandomUserAgent(c)

	if args.BadProxiesOutputFile != nil {
		c.OnError(func(response *colly.Response, err error) {
			log.Printf("[DEBUG] Bad proxy found error status %d [%s] [%s]", response.StatusCode, response.Request.ProxyURL, err)
			args.writeBadProxy(response.Request.ProxyURL)
		})
	}

	c.OnResponse(func(response *colly.Response) {
		log.Printf("[DEBUG] Good proxy found [%s]", response.Request.ProxyURL)
		args.writeGoodProxy(response.Request.ProxyURL)
	})

	for _, proxy := range args.ProxyList {
		log.Printf("[DEBUG] Doing %s", proxy)
		if err = c.Visit(args.WebsiteToCrawl); err != nil {
			log.Printf("Error happening doing Visit %s", err)
		}
	}

	// Wait until threads are finished
	c.Wait()
}
