package main

import (
	// "encoding/json"
	"fmt"
	"net/http"
	// import the data package which contains the definition for Comment
	// _ "github.com/georgie5/API-comments/internal/data"
)

func (a *applicationDependecies) createCommentHandler(w http.ResponseWriter,
	r *http.Request) {
	// create a struct to hold a comment
	// we use struct tags[“] to make the names display in lowercase
	var incomingData struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	}
	// perform the decoding
	// err := json.NewDecoder(r.Body).Decode(&incomingData)
	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// for now display the result
	fmt.Fprintf(w, "%+v\n", incomingData)
}
