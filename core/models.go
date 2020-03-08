package core

import (
	"time"
)

type ProxyChecker struct {
	ProxyList    []string
	TimeoutProxy time.Duration
	AvailableCheckers
	WebsiteCloudFlare string
}

type AvailableCheckers struct {
	CloudFlareBypass bool
	Basic            bool
}
