package sample_component

import (
	"fmt"
	"net/http"
)

func NewDummyHandlerFunc(val string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"CUSTOM_VALUE": %q}`, val)
	}
}
