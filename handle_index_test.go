package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
				"X-Forwarded-For": "192.0.2.1:1234",
			},
			respCode: http.StatusOK,
			respBody: "ip: 192.0.2.1:1234\n",
		},
		{
			reqHeaders: map[string]string{
				"X-Forwarded-For": "192.0.2.1:1234",
				"User-Agent":      "curl/7.77.0",
			},
			respCode: http.StatusOK,
			respBody: "ip: 192.0.2.1:1234\nuser_agent: curl/7.77.0\n",
		},
	}

	mockGeoLite2 := &mockGeoLite2{}
	srv := server{
		router: http.NewServeMux(),
		dbASN:  mockGeoLite2,
		dbCity: mockGeoLite2,
	}
	srv.AddMiddleware(RealIPMiddleware)
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
