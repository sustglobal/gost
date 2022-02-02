package httpapi

import (
	"expvar"
	"net/http"

	"github.com/gorilla/mux"
)

func NewMetricsHandler() HandlerMounter {
	return &metricsHandler{
		handler: expvar.Handler(),
	}
}

type metricsHandler struct {
	handler http.Handler
}

func (h *metricsHandler) Mount(r *mux.Router) {
	r.Handle("/debug/vars", h.handler)
}

func (h *metricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP(w, r)
}
