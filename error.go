package api2go

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

// HTTPError is used for errors
type HTTPError struct {
	Err    error
	Msg    string
	Status int
	Errors []Error `json:"errors,omitempty"`
}

// Error can be used for all kind of application errors
// e.g. you would use it to define form errors or any
// other semantical application problems
// for more information see http://jsonapi.org/format/#errors
type Error struct {
	ID     string       `json:"id,omitempty"`
	Links  *ErrorLinks  `json:"links,omitempty"`
	Status string       `json:"Status,omitempty"`
	Code   string       `json:"code,omitempty"`
	Title  string       `json:"title,omitempty"`
	Detail string       `json:"detail,omitempty"`
	Source *ErrorSource `json:"source,omitempty"`
	Meta   interface{}  `json:"meta,omitempty"`
}

// ErrorLinks is used to provide an About URL that leads to
// further details about the particular occurrence of the problem.
//
// for more information see http://jsonapi.org/format/#error-objects
type ErrorLinks struct {
	About string `json:"about,omitempty"`
}

// ErrorSource is used to provide references to the source of an error.
//
// The Pointer is a JSON Pointer to the associated entity in the request
// document.
// The Paramter is a string indicating which query parameter caused the error.
//
// for more information see http://jsonapi.org/format/#error-objects
type ErrorSource struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

// marshalHTTPError marshals an internal httpError
func marshalHTTPError(input HTTPError) string {
	if len(input.Errors) == 0 {
		input.Errors = []Error{{Title: input.Msg, Status: strconv.Itoa(input.Status)}}
	}

	data, err := json.Marshal(input)

	if err != nil {
		log.Println(err)
		return "{}"
	}

	return string(data)
}

// NewHTTPError creates a new error with message and Status code.
// `Err` will be logged (but never sent to a client), `Msg` will be sent and `Status` is the http Status code.
// `Err` can be nil.
func NewHTTPError(err error, msg string, status int) HTTPError {
	return HTTPError{Err: err, Msg: msg, Status: status}
}

// Error returns a nice string represenation including the Status
func (e HTTPError) Error() string {
	msg := fmt.Sprintf("http error (%d) %s and %d more errors", e.Status, e.Msg, len(e.Errors))
	if e.Err != nil {
		msg += ", " + e.Err.Error()
	}

	return msg
}
