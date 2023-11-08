package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Authentication Successfully done. You can now get back to the CLI Terminal.")

	vars := mux.Vars(r)
	state := vars["state"]

	// Now you have access to the "state" query parameter
	fmt.Printf("User state: %s\n", state)
}
