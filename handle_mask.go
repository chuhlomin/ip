package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (s *server) handleMask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ipParam := chi.URLParam(r, "ip")
		maskParam := chi.URLParam(r, "mask")

		m, err := strconv.Atoi(maskParam)
		if err != nil {
			msg := fmt.Sprintf("mask is not int: %s", maskParam)
			log.Print(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		if m < 0 || m > 32 {
			msg := fmt.Sprintf("invalid mask: %s", maskParam)
			log.Print(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		_, ipv4Net, err := net.ParseCIDR(ipParam + "/" + maskParam)
		if err != nil {
			msg := fmt.Sprintf("invalid ip: %s", ipParam)
			log.Print(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, ipParam+"/"+maskParam+"\n\n")

		// print IP with paddings to align with binary format
		ipParts := strings.Split(ipParam, ".")
		for i, part := range ipParts {
			fmt.Fprintf(w, "%8s", part)
			if i < len(ipParts)-1 {
				fmt.Fprint(w, ".")
			}
		}
		fmt.Fprint(w, "\n")

		// print IP in binary format
		for i, part := range ipParts {
			p, err := strconv.Atoi(part)
			if err != nil {
				msg := fmt.Sprintf("invalid ip: %s", ipParam)
				log.Print(msg)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}

			binPart := fmt.Sprintf("%08b", p)
			fmt.Fprint(w, binPart)
			if i < len(ipParts)-1 {
				fmt.Fprint(w, ".")
			}
		}
		fmt.Fprint(w, "\n")

		// print mask
		for i := 0; i < m; i++ {
			fmt.Fprint(w, "X")
			if (i+1)%8 == 0 && i < m-1 {
				fmt.Fprint(w, ".")
			}
		}
		fmt.Fprint(w, "\n\n")

		if m >= 16 {
			// print all possible IPs
			mask := binary.BigEndian.Uint32(ipv4Net.Mask)
			start := binary.BigEndian.Uint32(ipv4Net.IP)

			// find the final address
			finish := (start & mask) | (mask ^ 0xffffffff)

			// loop through addresses as uint32
			for i := start; i <= finish; i++ {
				// convert back to net.IP
				ip := make(net.IP, 4)
				binary.BigEndian.PutUint32(ip, i)
				fmt.Fprintln(w, ip)
			}
		}
	}
}
