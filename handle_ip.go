package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"gopkg.in/yaml.v2"
)

type geoip struct {
	Country     string  `json:"country,omitempty" yaml:"country,omitempty"`
	CountryName string  `json:"country_name,omitempty" yaml:"country_name,omitempty"`
	City        string  `json:"city,omitempty" yaml:"city,omitempty"`
	Lat         float64 `json:"lat,omitempty" yaml:"lat,omitempty"`
	Lon         float64 `json:"lon,omitempty" yaml:"lon,omitempty"`
}

type asn struct {
	Number       uint   `json:"number,omitempty" yaml:"number,omitempty"`
	Organization string `json:"org,omitempty" yaml:"org,omitempty"`
}

type response struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent,omitempty" yaml:"user_agent,omitempty"`
	ASN       *asn   `json:"asn,omitempty" yaml:"asn,omitempty"`
	GeoIP     *geoip `json:"geoip,omitempty" yaml:"geoip,omitempty"`
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

		resp := &response{
			IP:        ip,
			UserAgent: r.UserAgent(),
		}

		nip := net.ParseIP(ip)

		resp.ASN = s.getASN(nip)
		resp.GeoIP = s.getCity(nip, lang)

		switch {
		case format == "json", r.Header.Get("Accept") == "application/json":
			b, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(b))
		case format == "yaml":
			b, err := yaml.Marshal(resp)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprint(w, string(b))
		}
	}
}

func (s *server) getASN(ip net.IP) *asn {
	if ip == nil {
		return nil
	}

	if s.dbASN == nil {
		return nil
	}

	a, err := s.dbASN.ASN(ip)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &asn{
		Number:       a.AutonomousSystemNumber,
		Organization: a.AutonomousSystemOrganization,
	}
}

func (s *server) getCity(ip net.IP, lang string) *geoip {
	if ip == nil {
		return nil
	}

	if s.dbCity == nil {
		return nil
	}

	city, err := s.dbCity.City(ip)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &geoip{
		Country:     city.Country.IsoCode,
		CountryName: city.Country.Names[lang],
		City:        city.City.Names[lang],
		Lat:         city.Location.Latitude,
		Lon:         city.Location.Longitude,
	}
}
