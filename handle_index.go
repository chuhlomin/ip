package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"gopkg.in/yaml.v2"
)

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

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, string(b))
	}
}
