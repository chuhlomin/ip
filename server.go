package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"runtime/debug"

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
	buildRevision, buildTime := buildInfo()

	s.router.HandleFunc("/", s.handleIndex())
	s.router.HandleFunc("/whois", s.handleWhois())
	s.router.HandleFunc("/{ip:[0-9.]+}json", s.handleIP("json"))
	s.router.HandleFunc("/{ip:[0-9.]+}", s.handleIP("yaml"))
	s.router.HandleFunc("/{ip:[0-9.]+}/whois", s.handleIPWhois())
	s.router.HandleFunc("/{ip:[0-9.]+}/{mask:[0-9]+}", s.handleMask())
	s.router.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/help")
		w.WriteHeader(http.StatusFound)
	})
	s.router.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/help")
		w.WriteHeader(http.StatusFound)
	})
	s.router.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/help")
		w.WriteHeader(http.StatusFound)
	})
	s.router.HandleFunc("/help", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, `ip.chuhlomin.com is a service for finding information about IP addresses.

It uses:
  * GeoLite2 databases for ASN and GeoIP lookups,
  * whois.iana.org for Whois queries.

Available endpoints:
  /help - this page
  / - index page, redirects to /{ip}, where {ip} is your IP address
  /{ip} - returns information about the IP address: ASN and GeoIP
  /whois - redirects to /{ip}/whois if IP is known, otherwise returns 404
  /{ip}/whois - returns the Whois information for the IP address
  /{ip}/{mask} - displays the IP in binary format, visualizing the mask

Example usage:
  curl -L https://ip.chuhlomin.com/
  curl https://ip.chuhlomin.com/1.1.1.1
  curl https://ip.chuhlomin.com/1.1.1.1/whois
  curl https://ip.chuhlomin.com/192.168.0.0/24

Version: 1.0.0
  Revision: %s
  Build time: %s
Source code: https://github.com/chuhlomin/ip
Author: Konstantin Chukhlomin
License: MIT

---

Known alternatives:
  https://ip.me
  https://ifconfig.co
  https://httpbin.org/ip
  https://ipinfo.io
  https://whatismyipaddress.com
`,
			buildRevision, buildTime)
	})
	s.router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "User-agent: *\nDisallow: /\n")
	})
	s.router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./favicon.ico")
	})
	s.router.HandleFunc("/og.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./og.png")
	})

	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		// check if the request path is a valid IP address
		ip := net.ParseIP(strings.Trim(r.URL.Path[1:], " /"))
		if ip == nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// redirect to the IP address page
		w.Header().Set("Location", "/"+ip.String())
		w.WriteHeader(http.StatusFound)
	})
}

func buildInfo() (revision, time string) {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
			case "vcs.time":
				time = setting.Value
			}
		}
	}

	return
}
