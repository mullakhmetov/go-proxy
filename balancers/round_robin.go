package balancers

import (
	"container/ring"
	"net/url"
	"sync"
)

type RoundRobin struct {
	ring *ring.Ring
	l    sync.Mutex
}

func NewRoundRobin(urls_str []string) *RoundRobin {
	r := ring.New(len(urls_str))
	for _, v := range urls_str {
		url, err := url.Parse(v)
		if err != nil {
			panic(err.Error())
		}
		r.Value = Backend{Url: *url, Active: true}
		r = r.Next()
	}
	return &RoundRobin{ring: r}
}

func (rr *RoundRobin) ChooseTarget() Backend {
	rr.l.Lock()
	defer rr.l.Unlock()
	n := rr.ring.Value.(Backend)
	rr.ring = rr.ring.Next()
	return n
}

func (rr *RoundRobin) TargetsCount() int {
	return rr.ring.Len()
}
