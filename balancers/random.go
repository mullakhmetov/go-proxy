package balancers

import (
	"math/rand"
	"net/url"
)

type RandomBalancer struct {
	Backends []Backend
}

func NewRandomBalancer(urls_str []string) *RandomBalancer {
	backends := make([]Backend, len(urls_str))
	for i, v := range urls_str {
		url, err := url.Parse(v)
		if err != nil {
			panic(err.Error())
		}
		backends[i] = Backend{Url: *url, Active: true}
	}
	return &RandomBalancer{backends}
}

func (rb *RandomBalancer) ChooseTarget() Backend {
	return rb.Backends[rand.Int()%len(rb.Backends)]
}

func (rb *RandomBalancer) TargetsCount() int {
	return len(rb.Backends)
}
