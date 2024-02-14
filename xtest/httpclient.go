package th

import (
	"net/http"
)

// RoundTripFunc defines a function that wil be used on Transport of the http.Client.
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip implement the http.RoundTripper interface.
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestHTTPClient creates a *http.Client with the provided RoundTripFunc.
func NewTestHTTPClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}
