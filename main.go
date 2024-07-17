package main

import (
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/oschwald/geoip2-golang"
)

type config struct {
	Port             string        `env:"PORT" envDefault:"8080"`
	GeoLite2ASNPath  string        `env:"GEOLITE_ASN_PATH" envDefault:"GeoLite2/GeoLite2-ASN.mmdb"`
	GeoLite2CityPath string        `env:"GEOLITE_CITY_PATH" envDefault:"GeoLite2/GeoLite2-City.mmdb"`
	WhoisTimeout     time.Duration `env:"WHOIS_TIMEOUT" envDefault:"5s"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	dbASN, err := geoip2.Open(cfg.GeoLite2ASNPath)
	if err != nil {
		log.Println(err)
	}

	dbCity, err := geoip2.Open(cfg.GeoLite2CityPath)
	if err != nil {
		log.Println(err)
	}

	srv := server{
		router: http.NewServeMux(),
		dbASN:  dbASN,
		dbCity: dbCity,
		whois:  &WhoisClient{},
	}
	srv.AddMiddleware(RealIPMiddleware)
	srv.AddMiddleware(LoggerMiddleware)
	srv.routes()

	log.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, &srv); err != nil {
		log.Fatal(err)
	}
}
