package core

import (
	"github.com/kangoo13/proxy-checker/checkers"
)

func CheckProxies(proxyChecker *ProxyChecker) []string {
	var checkersWrapper = checkers.Wrapper{
		ProxiesToTest:   proxyChecker.ProxyList,
		AcceptedProxies: []string{},
		TimeoutProxy:    proxyChecker.TimeoutProxy,
	}

	if proxyChecker.Basic {
		checkersWrapper.WebsiteToCrawl = "https://httpbin.org/get"
		checkersWrapper.PrepareNextChecker()
		checkers.CheckIsAlive(&checkersWrapper)
	}

	if proxyChecker.CloudFlareBypass {
		checkersWrapper.WebsiteToCrawl = proxyChecker.WebsiteCloudFlare
		checkersWrapper.PrepareNextChecker()
		checkers.CheckCloudFlareBypass(&checkersWrapper)
	}

	return checkersWrapper.AcceptedProxies
}
