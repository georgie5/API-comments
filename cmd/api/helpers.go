package main

import (
	"encoding/json"
	"net/http"
)

func (a *applicationDependecies) writeJSON(w http.ResponseWriter,
	status int, data any, headers http.Header) error {

	jsResponse, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	jsResponse = append(jsResponse, '\n')
	//additional headers to be set
	for key, value := range headers {
		w.Header()[key] = value
	}

	//set content type header
	w.Header().Set("Content-Type", "application/json")
	//explicitly set the respnonse status code
	w.WriteHeader(status)
	_, err = w.Write(jsResponse)
	if err != nil {
		return err
	}

	return nil

}
