package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/require"
)

func TestHandlerIndex(t *testing.T) {
	tests := []struct {
		reqHeaders  map[string]string
		respCode    int
		respBody    string
		respHeaders map[string]string
	}{
		{
			reqHeaders: map[string]string{
				"X-Forwarded-For": "1.2.3.4",
			},
			respCode: http.StatusFound,
			respHeaders: map[string]string{
				"Location": "/1.2.3.4",
			},
		},
		{
			reqHeaders: map[string]string{
				"X-Forwarded-For": "[::1]:62993",
			},
			respCode: http.StatusOK,
			respBody: "ip: '[::1]:62993'\n",
		},
		{
			reqHeaders: map[string]string{
				"X-Forwarded-For": "[::1]:62993",
				"User-Agent":      "curl/7.77.0",
			},
			respCode: http.StatusOK,
			respBody: "ip: '[::1]:62993'\nuser_agent: curl/7.77.0\n",
		},
	}

	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	mockGeoLite2 := &mockGeoLite2{}
	srv := server{
		router: router,
		dbASN:  mockGeoLite2,
		dbCity: mockGeoLite2,
	}
	srv.routes()

	for _, tt := range tests {
		req := httptest.NewRequest("GET", "/", nil)
		for k, v := range tt.reqHeaders {
			req.Header.Set(k, v)
		}
		w := httptest.NewRecorder()

		srv.ServeHTTP(w, req)

		require.Equal(t, tt.respCode, w.Code)
		for k, v := range tt.respHeaders {
			require.Equal(t, v, w.Header().Get(k))
		}
		require.Equal(t, tt.respBody, w.Body.String())

		mockGeoLite2.AssertNotCalled(t, "ASN")
		mockGeoLite2.AssertNotCalled(t, "City")
	}
}
