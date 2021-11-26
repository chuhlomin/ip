package main

import (
	"net"

	"github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/mock"
)

type mockGeoLite2 struct {
	mock.Mock
}

func (mg *mockGeoLite2) City(ip net.IP) (*geoip2.City, error) {
	args := mg.Called(ip)
	return args.Get(0).(*geoip2.City), args.Error(1)
}

func (mg *mockGeoLite2) ASN(ip net.IP) (*geoip2.ASN, error) {
	args := mg.Called(ip)
	return args.Get(0).(*geoip2.ASN), args.Error(1)
}
