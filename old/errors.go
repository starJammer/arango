package old

import (
	"encoding/json"
)

//ArangoError is the base set of json fields in a typical arango response.
//All the api functions will return an ArangoError when something
//bad happens. Nil will be returned when things went as planned.
//If an error happened BEFORE we consulted the REST API, meaning
//the error happened in my code or because of parameter checks
//before the http request, then
//Code and ErrorNum will be -1. Otherwise, if we succeeded in making
//the request to the server, they will contain whatever the server/api
//would normally return.
type ArangoError struct {
	IsError      bool   `json:"error"`
	Code         int    `json:"code"`
	ErrorNum     int    `json:"errorNum"`
	ErrorMessage string `json:"errorMessage"`

	//Reserved for some operations that will return
	//the id, rev, or key of the document.
	//For example, /_api/document/{doc-handle} when it
	//return a 412 error
	Id  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	Key string `json:"_key,omitempty"`
}

func (a ArangoError) Error() string {
	b, _ := json.Marshal(a)
	return string(b)
}

func newError(msg string) ArangoError {
	return ArangoError{
		IsError:      true,
		Code:         -1,
		ErrorNum:     -1,
		ErrorMessage: msg,
	}
}
