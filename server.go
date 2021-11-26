package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oschwald/geoip2-golang"
)

type GeoIPReader interface {
	City(net.IP) (*geoip2.City, error)
	ASN(net.IP) (*geoip2.ASN, error)
}

type Whois interface {
	Query(string) (string, error)
}

type server struct {
	router chi.Router
	dbASN  GeoIPReader
	dbCity GeoIPReader
	whois  Whois
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) routes() {
	s.router.HandleFunc("/", s.handleIndex())
	s.router.HandleFunc("/{ip:[0-9.]+}json", s.handleIP("json"))
	s.router.HandleFunc("/{ip:[0-9.]+}", s.handleIP("yaml"))
	s.router.HandleFunc("/{ip:[0-9.]+}/whois", s.handleWhois())
	s.router.HandleFunc("/{ip:[0-9.]+}/{mask:[0-9]+}", s.handleMask())
	s.router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "User-agent: *\nDisallow: /\n")
	})
}
