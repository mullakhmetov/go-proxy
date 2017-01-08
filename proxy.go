package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/mullakhmetov/goxy/balancers"
)

type ProxyHandler struct {
	config   ServiceConfig
	balancer balancers.Balancer
}

func NewMultipleReverseProxy(h *ProxyHandler) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		target_url := h.balancer.ChooseTarget().Url
		req.URL.Scheme = target_url.Scheme
		req.URL.Host = target_url.Host
		fmt.Println("Call: Director.", target_url.String())
	}
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			fmt.Println("Call: Transport.Proxy")
			return http.ProxyFromEnvironment(req)
		},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			fmt.Println("Call: Transport.DialContext")
			attempts := h.balancer.TargetsCount()
			var conn net.Conn
			var err error
			for i := 0; ; i++ {
				conn, err = (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext(ctx, network, addr)
				if err == nil {
					fmt.Println("Error during DIAL:", err.Error())
					return conn, nil
				}
				if i >= (attempts - 1) {
					break
				}
				target_url := h.balancer.ChooseTarget().Url
				addr = target_url.Host + target_url.Path
				fmt.Println("Retrying with...", addr)
			}
			return conn, err
		},
	}

	return &httputil.ReverseProxy{
		Director:  director,
		Transport: transport,
	}
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Call: ServeHTTP")
	proxy := NewMultipleReverseProxy(h)
	proxy.ServeHTTP(w, r)
}

func StartProxy(config *Config) {
	wg := &sync.WaitGroup{}
	for _, v := range config.Services {
		wg.Add(1)
		port := fmt.Sprintf(":%s", v.Bind)
		balancer, err := balancers.CreateBalancer(v.Balance, v.Backends)
		if err != nil {
			log.Fatal(err.Error())
		}
		handler := &ProxyHandler{config: v, balancer: balancer}
		go func() {
			defer wg.Done()
			fmt.Println("binding", port)
			log.Fatal(http.ListenAndServe(port, handler))
		}()
	}
	wg.Wait()
}
