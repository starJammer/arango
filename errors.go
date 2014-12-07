package arango

import (
    "encoding/json"
)

//ArangoError is the base set of fields in a typical arango response
//All the api functions will return an ArangoError when something
//bad happens. Nil will be returned when things went as planned.
type ArangoError struct {
	IsError      bool `json:"error"`
    Code         int  `json:"code"`
    ErrorNum     int   `json:"errorNum"`
    ErrorMessage string `json:"errorMessage"`
}

func (a ArangoError) Error() string {
    b, _ := json.Marshal( a )
	return string( b )
}

func newError( msg string ) ArangoError{
    return ArangoError{
        IsError : true,
        Code : -1,
        ErrorNum : -1,
        ErrorMessage : msg,
    }
}
