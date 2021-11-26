package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestHandlerMask(t *testing.T) {
	tests := []struct {
		path string
		resp string
	}{
		{
			path: "/1.2.3.4/32",
			resp: `1.2.3.4/32

       1.       2.       3.       4
00000001.00000010.00000011.00000100
XXXXXXXX.XXXXXXXX.XXXXXXXX.XXXXXXXX

1.2.3.4
`,
		},
		{
			path: "/1.2.3.4/30",
			resp: `1.2.3.4/30

       1.       2.       3.       4
00000001.00000010.00000011.00000100
XXXXXXXX.XXXXXXXX.XXXXXXXX.XXXXXX

1.2.3.4
1.2.3.5
1.2.3.6
1.2.3.7
`,
		},
	}

	srv := server{router: chi.NewRouter()}
	srv.routes()

	for _, tt := range tests {
		req := httptest.NewRequest("GET", tt.path, nil)
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, tt.resp, w.Body.String())
	}
}

func TestHandlerMaskInvalidMask(t *testing.T) {
	srv := server{router: chi.NewRouter()}
	srv.routes()

	req := httptest.NewRequest("GET", "/1.2.3.4/33", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Equal(t, "invalid mask: 33\n", w.Body.String())
}

func TestHandlerMaskInvalidCIDR(t *testing.T) {
	srv := server{router: chi.NewRouter()}
	srv.routes()

	req := httptest.NewRequest("GET", "/256.0.0.0/32", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Equal(t, "invalid ip: 256.0.0.0\n", w.Body.String())
}
