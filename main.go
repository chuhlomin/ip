package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/oschwald/geoip2-golang"
	"gopkg.in/yaml.v2"
)

type config struct {
	GeoLite2ASNPath  string        `env:"GEOLITE_ASN_PATH" envDefault:"GeoLite2/GeoLite2-ASN.mmdb"`
	GeoLite2CityPath string        `env:"GEOLITE_CITY_PATH" envDefault:"GeoLite2/GeoLite2-City.mmdb"`
	WhoisTimeout     time.Duration `env:"WHOIS_TIMEOUT" envDefault:"5s"`
}

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

func getIP(path string) net.IP {
	if ip := path; ip != "/" {
		if ip = ip[1:]; len(ip) > 0 {
			if netIP := net.ParseIP(ip); netIP != nil {
				return netIP
			}
		}
	}
	return nil
}

func digIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if len(ip) == 0 {
		return r.RemoteAddr
	}
	return ip
}

func buildResponse(ip net.IP, r *http.Request, lang string) (resp *response) {
	resp = &response{
		IP:        ip.String(),
		UserAgent: r.UserAgent(),
	}

	if ip != nil {
		asn, err := dbASN.ASN(ip)
		if err != nil {
			log.Println(err)
			return
		}
		resp.ASN.Number = asn.AutonomousSystemNumber
		resp.ASN.Organization = asn.AutonomousSystemOrganization

		city, err := dbCity.City(ip)
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

func handler(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r.URL.Path)
	if ip == nil {
		ip = net.ParseIP(digIP(r))
		if ip != nil {
			// redirect to IP
			w.Header().Set("Location", "/"+ip.String())
			w.WriteHeader(http.StatusFound)
			return
		}
	}

	lang := "en"
	if r.URL.Query().Get("lang") != "" {
		lang = r.URL.Query().Get("lang")
	}

	resp := buildResponse(ip, r, lang)

	b, err := yaml.Marshal(resp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/yaml")
	fmt.Fprint(w, string(b))
}

var (
	dbASN  *geoip2.Reader
	dbCity *geoip2.Reader
)

func main() {
	var err error
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	dbASN, err = geoip2.Open(cfg.GeoLite2ASNPath)
	if err != nil {
		log.Fatal(err)
	}

	dbCity, err = geoip2.Open(cfg.GeoLite2CityPath)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handler)

	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "User-agent: *\nDisallow: /\n")
	})

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
