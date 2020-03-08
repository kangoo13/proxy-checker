package checkers

import (
	"github.com/gocolly/colly"
	"log"
	"net/url"
)

func CheckCloudFlareBypass(args *Wrapper) {
	var onErrorCallback = func(response *colly.Response, err error) {
		log.Printf("[DEBUG CheckCloudFlareBypass] Bad proxy found error status %d [%s] [%s]", response.StatusCode, response.Request.ProxyURL, err)
	}

	var onResponseCallback = func(response *colly.Response) {
		log.Printf("[DEBUG CheckCloudFlareBypass] Good proxy found [%s]", response.Request.ProxyURL)
		args.AddGoodProxy(response.Request.ProxyURL)
	}

	var collectors []*colly.Collector

	for _, proxy := range args.ProxiesToTest {
		newCollector := NewCollector(args)
		newCollector.OnResponse(onResponseCallback)
		newCollector.OnError(onErrorCallback)
		u, err := url.Parse(proxy)
		if err != nil {
			log.Printf("[CheckCloudFlareBypass] error while parseing proxy %s %s", proxy, err)
			newCollector = nil
			continue
		}
		newCollector.SetProxyFunc(ProxyURL(u))
		if err = newCollector.Visit(args.WebsiteToCrawl); err != nil {
			log.Printf("[CheckCloudFlareBypass] error happening doing Visit %s", err)
			newCollector = nil
			continue

		}
		collectors = append(collectors, newCollector)
	}

	for _, collec := range collectors {
		collec.Wait()
	}

}
