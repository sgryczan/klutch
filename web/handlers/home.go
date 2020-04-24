package handlers

import "net/http"

// Homehandler redirects to swaggerui
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api/", 302)
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "This is home")
}
