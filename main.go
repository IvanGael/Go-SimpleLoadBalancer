package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type LoadBalancer struct {
	backendUrls []*url.URL
	currentIdx  int
	mux         sync.Mutex
}

func NewLoadBalancer(urls []string) *LoadBalancer {
	lb := &LoadBalancer{}
	for _, u := range urls {
		parsedUrl, err := url.Parse(u)
		if err != nil {
			log.Fatalf("Error parsing URL: %s", err)
		}
		lb.backendUrls = append(lb.backendUrls, parsedUrl)
	}
	return lb
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lb.mux.Lock()
	defer lb.mux.Unlock()

	// Choose backend server using round-robin algorithm
	backendUrl := lb.backendUrls[lb.currentIdx]
	lb.currentIdx = (lb.currentIdx + 1) % len(lb.backendUrls)

	// Proxy the request to the chosen backend server
	reverseProxy := &httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = backendUrl.Scheme
		req.URL.Host = backendUrl.Host
		req.URL.Path = backendUrl.Path
	}}
	reverseProxy.ServeHTTP(w, r)
}

func main() {
	backendUrls := []string{
		"http://localhost:8001",
		"http://localhost:8002",
	}

	loadBalancer := NewLoadBalancer(backendUrls)

	log.Println("Load Balancer started on :8080")
	log.Fatal(http.ListenAndServe(":8080", loadBalancer))
}
