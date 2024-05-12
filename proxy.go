package roxy

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRequestID(length int) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func requestIdFromHeaders(set http.Header) string {

	if requestID := set.Get("X-B3-TraceId"); requestID != "" {
		return requestID
	}
	return generateRequestID(8)
}

type Roxy struct {
	cfg *Config
}

func (x *Roxy) accessDenied(w http.ResponseWriter, data *AccessDenied) {
	w.Header().Set("Content-Type", "text/html")
	accessDeniedTemplate.Execute(w, data)
}

func (x *Roxy) handler() http.HandlerFunc {
	targetURL, err := url.Parse(x.cfg.Target)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(w http.ResponseWriter, r *http.Request) {
		// intercept health check call
		if r.URL.Path == x.cfg.HealthCheckPath {
			io.WriteString(w, "OK")
			return
		}

		requestID := requestIdFromHeaders(r.Header)
		log.Printf("[%s] Start\n", requestID)
		log.Printf("[%s] Forwarded URL: %s", requestID, r.URL.Path)

		x.accessDenied(w, &AccessDenied{EmailName: x.cfg.EmailName, Email: x.cfg.Email, ForwardedURL: r.URL.Path, RequestID: requestID, ClientID: r.RemoteAddr})
		return

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

}

func (x *Roxy) Start() {
	// Start the server on port 8080
	// "handler" is now the proxy
	http.HandleFunc("/", x.handler())
	address := fmt.Sprintf("%s:%s", x.cfg.Host, x.cfg.Port)
	log.Printf("Starting proxy server on %s", address)
	err := http.ListenAndServe(fmt.Sprintf(address), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func NewRoxy(cfg *Config) *Roxy {
	return &Roxy{cfg: cfg}
}
