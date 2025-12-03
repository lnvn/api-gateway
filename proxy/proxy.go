package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewProxy creates a new reverse proxy for a given target URL.
func NewProxy(target string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}

// ProxyHandler returns a http.HandlerFunc that forwards requests to the proxy.
func ProxyHandler(proxy *httputil.ReverseProxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Update the headers to allow for SSL redirection
		r.Host = r.URL.Host
		proxy.ServeHTTP(w, r)
	}
}
