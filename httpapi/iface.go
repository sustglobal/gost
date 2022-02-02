package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
)

type HandlerMounter interface {
	http.Handler
	Mount(*mux.Router)
}
