package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/oschwald/geoip2-golang"
)

type config struct {
	GeoLite2ASNPath  string `env:"GEOLITE_ASN_PATH" envDefault:"GeoLite2/GeoLite2-ASN.mmdb"`
	GeoLite2CityPath string `env:"GEOLITE_CITY_PATH" envDefault:"GeoLite2/GeoLite2-City.mmdb"`
}

func getIP(r *http.Request) (netIP net.IP) {
	if ip := r.URL.Path; ip != "/" {
		if ip = ip[1:]; len(ip) > 0 {
			if netIP = net.ParseIP(ip); netIP != nil {
				return
			}
		}
	}

	ip := r.Header.Get("X-Forwarded-For")
	if len(ip) == 0 {
		ip = r.RemoteAddr
	}

	return net.ParseIP(ip)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/yaml")

	lang := "en"
	if r.URL.Query().Get("lang") != "" {
		lang = r.URL.Query().Get("lang")
	}

	ip := getIP(r)
	fmt.Fprintf(w, "IP: %v\n", ip.String())
	fmt.Fprintf(w, "User-Agent: %v\n", r.UserAgent())

	if ip != nil {
		asn, err := dbASN.ASN(ip)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Fprintf(w, "ASN: %v\n", asn.AutonomousSystemNumber)
		fmt.Fprintf(w, "ASN Name: %v\n", asn.AutonomousSystemOrganization)

		city, err := dbCity.City(ip)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Fprintf(w, "Country: %v\n", city.Country.IsoCode)
		fmt.Fprintf(w, "Country Name: %v\n", city.Country.Names[lang])
		fmt.Fprintf(w, "City: %v\n", city.City.Names[lang])
		fmt.Fprintf(w, "Latitude: %v\n", city.Location.Latitude)
		fmt.Fprintf(w, "Longitude: %v\n", city.Location.Longitude)
	}
}

var dbASN *geoip2.Reader
var dbCity *geoip2.Reader

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

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
