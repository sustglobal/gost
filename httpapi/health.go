package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewHealthHandler() HandlerMounter {
	return new(healthHandler)
}

type healthHandler struct{}

func (h *healthHandler) Mount(r *mux.Router) {
	r.Handle("/healthz", h)
}

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
