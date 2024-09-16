package main

import (
	"encoding/json"
	"net/http"
)

func (a *applicationDependecies) healthChechHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":      "available",
		"environment": a.config.environment,
		"version":     appVersion,
	}

	jsResponse, err := json.Marshal(data)

}
