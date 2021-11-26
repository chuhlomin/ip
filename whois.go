package main

import (
	"github.com/likexian/whois"
	"github.com/stretchr/testify/mock"
)

type WhoisClient struct{}

func (w *WhoisClient) Query(ip string) (string, error) {
	return whois.Whois(ip)
}

type mockWhois struct {
	mock.Mock
}

func (mw *mockWhois) Query(ip string) (string, error) {
	args := mw.Called(ip)
	return args.String(0), args.Error(1)
}
