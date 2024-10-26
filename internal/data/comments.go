package data

import (
	"time"

	"georgie5.net/API-comments/internal/validator"
)

// each name begins with uppercase so that they are exportable/public
type Comment struct {
	ID        int64     `json: "id"`      // unique value for each comment
	Content   string    `json: "content"` // the comment data
	Author    string    `json: "author"`  // the person who wrote the comment
	CreatedAt time.Time `json: "-"`       // database timestamp
	Version   int32     `json: "version"` // incremented on each update
}

// Create a function that performs the validation checks
func ValidateComment(v *validator.Validator, comment *Comment) {

	// check if the Content field is empty
	v.Check(comment.Content != "", "content", "must be provided")
	// check if the Author field is empty
	v.Check(comment.Author != "", "author", "must be provided")
	// check if the Content field is empty
	v.Check(len(comment.Content) <= 100, "content", "must not be more than 100 bytes long")
	// check if the Author field is empty
	v.Check(len(comment.Author) <= 25, "author", "must not be more than 25 bytes long")

}
