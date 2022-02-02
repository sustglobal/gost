package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestMetricsHandler(t *testing.T) {
	rtr := mux.NewRouter()

	h := NewMetricsHandler()
	h.Mount(rtr)

	req, err := http.NewRequest("GET", "/debug/vars", nil)
	if err != nil {
		t.Fatalf("Failed building HTTP request: %v", err)
	}

	rec := httptest.NewRecorder()
	rtr.ServeHTTP(rec, req)

	res := rec.Result()

	if res.StatusCode != 200 {
		t.Errorf("Received unexpected status code: want=200 got=%d", res.StatusCode)
	}
}
