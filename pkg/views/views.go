package views

import (
	"fmt"
	"net/http"
)

func Hello(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintln(w, "Hello World!")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
	}
}
