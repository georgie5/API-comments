package main

import (
	"fmt"
	"net/http"
)

func (a *applicationDependecies) healthChechHandler(w http.ResponseWriter, r *http.Request) {

	// fmt.Fprintln(w, "status: available")
	// fmt.Fprintf(w, "environment: %s\n", a.config.environmnet)
	// fmt.Fprintf(w, "version: %s\n", appVersion)

	jsResponse := `{"status": "available", "environment": %q, "version": %q}`
	jsResponse = fmt.Sprintf(jsResponse, a.config.environmnet, appVersion)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsResponse))

}
