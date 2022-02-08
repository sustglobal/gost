package sample_component

import "net/http"

func DummyHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}
