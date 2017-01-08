package balancers

import (
	"errors"
	"fmt"
	"net/url"
)

type Balancer interface {
	ChooseTarget() Backend
	TargetsCount() int
}

type Backend struct {
	Url    url.URL
	Active bool
}

func CreateBalancer(t string, urls []string) (Balancer, error) {
	switch t {
	case "roundrobin":
		return NewRoundRobin(urls), nil
	case "random":
		return NewRandomBalancer(urls), nil
	default:
		return nil, errors.New(fmt.Sprintf("Invalid balancer type: %s", t))
	}
}
