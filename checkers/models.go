package checkers

import (
	"sync"
	"time"
)

type Wrapper struct {
	WebsiteToCrawl  string
	ProxiesToTest   []string
	AcceptedProxies []string
	TimeoutProxy    time.Duration
	mu              sync.Mutex
	checks          int
}

func (a *Wrapper) AddGoodProxy(proxy string) {
	a.mu.Lock()

	a.AcceptedProxies = append(a.AcceptedProxies, proxy)

	a.mu.Unlock()
}

func (a *Wrapper) PrepareNextChecker() {
	if a.checks != 0 {
		a.ProxiesToTest = a.AcceptedProxies
		a.AcceptedProxies = []string{}
	}
	a.checks += 1
}
