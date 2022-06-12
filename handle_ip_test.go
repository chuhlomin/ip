package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerIP(t *testing.T) {
	city := &geoip2.City{}
	city.Country.IsoCode = "US"
	city.Country.Names = map[string]string{"en": "United States"}
	city.City.Names = map[string]string{"en": "City"}
	city.Location.Latitude = 1.0
	city.Location.Longitude = 2.0

	mockGeoLite2 := &mockGeoLite2{}
	mockGeoLite2.On("ASN", mock.Anything).Return(&geoip2.ASN{
		AutonomousSystemNumber:       123,
		AutonomousSystemOrganization: "AS Organization",
	}, nil)
	mockGeoLite2.On("City", mock.Anything).Return(city, nil)

	srv := server{
		router: chi.NewRouter(),
		dbASN:  mockGeoLite2,
		dbCity: mockGeoLite2,
	}
	srv.routes()

	req := httptest.NewRequest("GET", "/1.2.3.4", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, `ip: 1.2.3.4
asn:
    number: 123
    org: AS Organization
geoip:
    country: US
    country_name: United States
    city: City
    lat: 1
    lon: 2
`, w.Body.String())
}

func TestHandlerIPNoGeoIPDatabases(t *testing.T) {
	srv := server{
		router: chi.NewRouter(),
		dbASN:  nil,
		dbCity: nil,
	}
	srv.routes()

	req := httptest.NewRequest("GET", "/1.2.3.4", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, `ip: 1.2.3.4
`, w.Body.String())
}

func TestHandlerIPJSON(t *testing.T) {
	srv := server{router: chi.NewRouter()}
	srv.routes()

	req := httptest.NewRequest("GET", "/1.2.3.4.json", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, `{
  "ip": "1.2.3.4"
}`, w.Body.String())
}

func TestHandlerIPJSONAccept(t *testing.T) {
	srv := server{router: chi.NewRouter()}
	srv.routes()

	req := httptest.NewRequest("GET", "/1.2.3.4", nil)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, `{
  "ip": "1.2.3.4"
}`, w.Body.String())
}

func TestHandlerIPGeoIPError(t *testing.T) {
	mockGeoLite2 := &mockGeoLite2{}
	mockGeoLite2.On("ASN", mock.Anything).Return(&geoip2.ASN{}, fmt.Errorf("ASN error"))
	mockGeoLite2.On("City", mock.Anything).Return(&geoip2.City{}, fmt.Errorf("City error"))

	srv := server{
		router: chi.NewRouter(),
		dbASN:  mockGeoLite2,
		dbCity: mockGeoLite2,
	}
	srv.routes()

	req := httptest.NewRequest("GET", "/1.2.3.4", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, `ip: 1.2.3.4
`, w.Body.String())
}

func TestHandlerIPIncorrectIP(t *testing.T) {
	mockGeoLite2 := &mockGeoLite2{}

	srv := server{
		router: chi.NewRouter(),
		dbASN:  mockGeoLite2,
		dbCity: mockGeoLite2,
	}
	srv.routes()

	req := httptest.NewRequest("GET", "/999.999.999.999", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, `ip: 999.999.999.999
`, w.Body.String())

	mockGeoLite2.AssertNotCalled(t, "ASN")
	mockGeoLite2.AssertNotCalled(t, "City")
}

func TestHandlerIPSupportLangParameter(t *testing.T) {
	city := &geoip2.City{}
	city.Country.IsoCode = "US"
	city.Country.Names = map[string]string{"en": "United States", "ru": "США"}
	city.City.Names = map[string]string{"en": "City"}
	city.Location.Latitude = 1.0
	city.Location.Longitude = 2.0

	mockGeoLite2 := &mockGeoLite2{}
	mockGeoLite2.On("ASN", mock.Anything).Return(&geoip2.ASN{
		AutonomousSystemNumber:       123,
		AutonomousSystemOrganization: "AS Organization",
	}, nil)
	mockGeoLite2.On("City", mock.Anything).Return(city, nil)

	srv := server{
		router: chi.NewRouter(),
		dbASN:  mockGeoLite2,
		dbCity: mockGeoLite2,
	}
	srv.routes()

	req := httptest.NewRequest("GET", "/1.2.3.4?lang=ru", nil)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, `ip: 1.2.3.4
asn:
    number: 123
    org: AS Organization
geoip:
    country: US
    country_name: США
    lat: 1
    lon: 2
`, w.Body.String())
}
