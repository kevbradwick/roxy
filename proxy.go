package roxy

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Roxy struct {
	cfg       *Config
	targetURL *url.URL
}

// generates a new request id, 8 char a-zA-Z0-9
func generateRequestID(length int) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// tries to get the request id from the header and if not present, generate a new one
func requestIdFromHeaders(set http.Header) string {
	if requestID := set.Get("X-B3-TraceId"); requestID != "" {
		return requestID
	}
	return generateRequestID(8)
}

// renders access denied response page
func (x *Roxy) accessDenied(requestID string, w http.ResponseWriter, r *http.Request) {
	templateData := &AccessDenied{
		ForwardedURL: r.URL.Path,
		RequestID:    requestID,
		ClientID:     r.RemoteAddr,
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusForbidden)
	accessDeniedTemplate.Execute(w, templateData)
}

// checks incoming request to see if it's coming from the load balancer health check
func (x *Roxy) isELBHealthChecker(r *http.Request, requestID string) bool {
	if v := r.Header.Get("User-Agent"); strings.HasPrefix(v, "ELB-HealthChecker") {
		log.Printf("[%s] Health check", requestID)
		return true
	}
	return false
}

// main handler. uses golang reverse proxy to process the request
func (x *Roxy) handler() http.HandlerFunc {

	proxy := httputil.NewSingleHostReverseProxy(x.targetURL)

	return func(w http.ResponseWriter, r *http.Request) {
		requestID := requestIdFromHeaders(r.Header)

		if x.isELBHealthChecker(r, requestID) {
			io.WriteString(w, "OK")
			return
		}

		log.Printf("[%s] Start\n", requestID)
		log.Printf("[%s] Forwarded URL: %s", requestID, r.URL.Path)

		forwardedFor := r.Header.Get("X-Forwarded-For")
		if forwardedFor == "" {
			r.RemoteAddr = "Unknown"
			log.Printf("[%s] missing header X-Forwarded-For\n", requestID)
			x.accessDenied(requestID, w, r)
			return
		}

		// ip address blocked
		if !IPAllowed(forwardedFor, x.cfg.AllowList) {
			log.Printf("[%s] Access denied for IP: %s", requestID, forwardedFor)
			x.accessDenied(requestID, w, r)
			return
		}

		// Update the headers to allow for SSL redirection
		r.URL.Host = x.targetURL.Host
		r.URL.Scheme = x.targetURL.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = x.targetURL.Host

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
	targetURL, err := url.Parse(cfg.Target)
	if err != nil {
		log.Fatal(err)
	}
	return &Roxy{cfg: cfg, targetURL: targetURL}
}
