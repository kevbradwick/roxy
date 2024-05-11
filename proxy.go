package roxy

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Roxy struct {
	cfg *Config
}

func (x *Roxy) handler(w http.ResponseWriter, r *http.Request) {

}

func (x *Roxy) Start() {
	targetURL, err := url.Parse(x.cfg.Target)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// A handler that will be called to handle the request
	handler := func(w http.ResponseWriter, r *http.Request) {

		// intercept health check call
		if r.URL.Path == x.cfg.HealthCheckPath {
			io.WriteString(w, "OK")
			return
		}

		// Update the headers to allow for SSL redirection
		r.URL.Host = targetURL.Host
		r.URL.Scheme = targetURL.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = targetURL.Host

		// Log the request; purely for debugging purposes to see it's going through the proxy
		log.Printf("Received request: %s %s %s", r.Method, r.Host, r.URL.Path)

		// ServeHttp is non blocking and uses a go routine under the hood
		proxy.ServeHTTP(w, r)
	}

	// Start the server on port 8080
	// "handler" is now the proxy
	http.HandleFunc("/", handler)
	log.Println("Starting proxy server on localhost:8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func NewRoxy(cfg *Config) *Roxy {
	return &Roxy{cfg: cfg}
}
