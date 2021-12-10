package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *server) handleWhois() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := net.ParseIP(r.RemoteAddr)
		if ip != nil {
			w.Header().Set("Location", "/"+ip.String()+"/whois")
			w.WriteHeader(http.StatusFound)
			return
		}
		http.NotFound(w, r)
		return
	}
}

func (s *server) handleIPWhois() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.whois == nil {
			log.Println("whois is nil")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		ip := chi.URLParam(r, "ip")
		whoisRaw, err := s.whois.Query(ip)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, whoisRaw)
	}
}
