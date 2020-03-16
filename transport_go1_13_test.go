// +build go1.13

package filtertransport

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func testFilteredRequst(t *testing.T, req *http.Request) {
	c := http.Client{Transport: DefaultTransport}
	resp, err := c.Do(req)

	if resp != nil {
		t.Fatal("expected resp to be nil for filtered ip")
	}
	var filterErr FilterError
	if !errors.As(err, &filterErr) {
		t.Fatal("expected err to be FilterError for filtered ip")
	}
}
func TestDefaultTransportWithContext(t *testing.T) {
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "http://127.0.0.1", nil)
	testFilteredRequst(t, req)
}

func TestDefaultTransportWithoutContext(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://127.0.0.1", nil)
	testFilteredRequst(t, req)
}

func TestTLSDefaultTransportWithContext(t *testing.T) {
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "https://127.0.0.1", nil)
	testFilteredRequst(t, req)
}

func TestTLSDefaultTransportWithoutContext(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://127.0.0.1", nil)
	testFilteredRequst(t, req)
}
