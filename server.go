package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/likexian/whois"
	"github.com/oschwald/geoip2-golang"
	"gopkg.in/yaml.v2"
)

type response struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent,omitempty" yaml:"user_agent,omitempty"`
	ASN       struct {
		Number       uint   `json:"number,omitempty" yaml:"number,omitempty"`
		Organization string `json:"org,omitempty" yaml:"org,omitempty"`
	} `json:"asn,omitempty" yaml:"asn,omitempty"`
	GeoIP struct {
		Country     string  `json:"country,omitempty" yaml:"country,omitempty"`
		CountryName string  `json:"country_name,omitempty" yaml:"country_name,omitempty"`
		City        string  `json:"city,omitempty" yaml:"city,omitempty"`
		Lat         float64 `json:"lat,omitempty" yaml:"lat,omitempty"`
		Lon         float64 `json:"lon,omitempty" yaml:"lon,omitempty"`
	} `json:"geoip,omitempty" yaml:"geoip,omitempty"`
}

type server struct {
	router chi.Router
	dbASN  *geoip2.Reader
	dbCity *geoip2.Reader
}

func (s *server) routes() {
	s.router.HandleFunc("/", s.handleIndex())
	s.router.HandleFunc("/{ip:[0-9.]+}json", s.handleIP("json"))
	s.router.HandleFunc("/{ip:[0-9.]+}", s.handleIP("yaml"))
	s.router.HandleFunc("/{ip:[0-9.]+}/whois", s.handleWhois())
	s.router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "User-agent: *\nDisallow: /\n")
	})
}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := net.ParseIP(r.RemoteAddr)
		if ip != nil {
			// redirect to IP
			w.Header().Set("Location", "/"+ip.String())
			w.WriteHeader(http.StatusFound)
			return
		}

		b, err := yaml.Marshal(response{
			IP:        r.RemoteAddr,
			UserAgent: r.UserAgent(),
		})
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/yaml")
		fmt.Fprint(w, string(b))
	}
}

func (s *server) handleIP(format string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := chi.URLParam(r, "ip")
		if format == "json" {
			// todo: fix regex
			ip = strings.TrimSuffix(ip, ".")
		}

		lang := "en"
		if r.URL.Query().Get("lang") != "" {
			lang = r.URL.Query().Get("lang")
		}

		resp := s.buildResponse(ip, r, lang)

		switch format {
		case "yaml":
			b, err := yaml.Marshal(resp)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/yaml")
			fmt.Fprint(w, string(b))

		case "json":
			b, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(b))
		}
	}
}

func (s *server) handleWhois() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := chi.URLParam(r, "ip")
		whoisRaw, err := whois.Whois(ip)
		if err != nil {
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, whoisRaw)
	}
}

func (s *server) buildResponse(ip string, r *http.Request, lang string) (resp *response) {
	resp = &response{
		IP:        ip,
		UserAgent: r.UserAgent(),
	}

	nip := net.ParseIP(ip)

	if nip != nil {
		asn, err := s.dbASN.ASN(nip)
		if err != nil {
			log.Println(err)
			return
		}
		resp.ASN.Number = asn.AutonomousSystemNumber
		resp.ASN.Organization = asn.AutonomousSystemOrganization

		city, err := s.dbCity.City(nip)
		if err != nil {
			log.Println(err)
			return
		}

		resp.GeoIP.Country = city.Country.IsoCode
		resp.GeoIP.CountryName = city.Country.Names[lang]
		resp.GeoIP.City = city.City.Names[lang]
		resp.GeoIP.Lat = city.Location.Latitude
		resp.GeoIP.Lon = city.Location.Longitude
	}

	return
}
