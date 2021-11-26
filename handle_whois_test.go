package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerWhois(t *testing.T) {
	mockWhois := &mockWhois{}
	mockWhois.On("Query", mock.Anything).Return(`WHOIS RAW RESPONSE`, nil)

	srv := server{
		router: chi.NewRouter(),
		whois:  mockWhois,
	}
	srv.routes()

	req := httptest.NewRequest("GET", "/1.2.3.4/whois", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, `WHOIS RAW RESPONSE`, w.Body.String())
	mockWhois.AssertCalled(t, "Query", "1.2.3.4")
}

func TestHandlerWhoisError(t *testing.T) {
	mockWhois := &mockWhois{}
	mockWhois.On("Query", mock.Anything).Return(``, fmt.Errorf("whois error"))

	srv := server{
		router: chi.NewRouter(),
		whois:  mockWhois,
	}
	srv.routes()

	req := httptest.NewRequest("GET", "/1.2.3.4/whois", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	mockWhois.AssertCalled(t, "Query", "1.2.3.4")
}

func TestHandlerWhoisNil(t *testing.T) {
	srv := server{
		router: chi.NewRouter(),
		whois:  nil,
	}
	srv.routes()

	req := httptest.NewRequest("GET", "/1.2.3.4/whois", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}
