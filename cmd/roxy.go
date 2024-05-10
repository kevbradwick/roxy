package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	fmt.Println("Roxy...")

	backendURL, err := url.Parse("https://www.bbc.co.uk") // Change this URL to your target backend server
	if err != nil {
		log.Fatal(err)
	}

	// Create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	// A handler that will be called to handle the request
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Update the headers to allow for SSL redirection
		r.URL.Host = backendURL.Host
		r.URL.Scheme = backendURL.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = backendURL.Host

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
