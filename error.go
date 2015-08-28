package arango

import (
	"fmt"
)

//ArangoError represents an arango error.
type ArangoError struct {
	//usually true because the api doesn't return an error
	//if this is false
	IsError bool `json:"error"`
	//The http response code
	Code int `json:"code"`
	//arango error number
	ErrorNum int `json:"errorNum"`
	//error description from arango server
	ErrorMessage string `json:"errorMessage"`

	//Reserved for some operations that will return
	//the id, rev, or key of the document in an error
	//json object.
	//This is primarily used when arango returns a 412
	//http respons code.
	//For example, GET /_api/document/{doc-handle} when it
	//return a 412 error
	Id  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
	Key string `json:"_key,omitempty"`
}

//Error is for the error interface in go
func (e ArangoError) Error() string {
	return fmt.Sprintf("isError: %t, code: %d, errorNum: %d, errorMessage: %s",
		e.IsError,
		e.Code,
		e.ErrorNum,
		e.ErrorMessage,
	)
}

func newArangoError(code int, message string) ArangoError {
	return ArangoError{
		IsError:      true,
		Code:         code,
		ErrorNum:     code,
		ErrorMessage: message,
	}
}
