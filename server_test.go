package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServerRobotsTXT(t *testing.T) {
	srv := server{router: http.NewServeMux()}
	srv.routes()

	req := httptest.NewRequest("GET", "/robots.txt", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, `User-agent: *
Disallow: /
`, w.Body.String())
}
