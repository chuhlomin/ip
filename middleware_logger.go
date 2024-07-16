package main

import (
	"log"
	"net/http"
	"time"
)

type Color string

const (
	red    Color = "\033[31m"
	green  Color = "\033[32m"
	yellow Color = "\033[33m"
	blue   Color = "\033[34m"
	cyan   Color = "\033[36m"
	reset  Color = "\033[0m"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}

type Entry struct {
	r    *http.Request
	time time.Time
}

func NewLogEntry(r *http.Request) *Entry {
	return &Entry{
		r:    r,
		time: time.Now(),
	}
}

func (e *Entry) Write(w *wrappedResponseWriter) {
	color := red
	switch {
	case w.StatusCode < 200:
		color = blue
	case w.StatusCode < 300:
		color = green
	case w.StatusCode < 400:
		color = cyan
	case w.StatusCode < 500:
		color = yellow
	}

	log.Printf("%s %s%d%s %s %s %s %v", e.r.RemoteAddr, color, w.StatusCode, reset, e.r.Method, e.r.URL.Path, e.r.Proto, time.Since(e.time))
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := &wrappedResponseWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		entry := NewLogEntry(r)
		defer entry.Write(wrapped)

		next.ServeHTTP(wrapped, r)
	})
}
