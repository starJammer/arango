package arango

import (
	"fmt"
)

type arangoError struct {
	IsErrorf      bool   `json:"error"`
	Codef         int    `json:"code"`
	ErrorNumf     int    `json:"errorNum"`
	ErrorMessagef string `json:"errorMessage"`
	Idf           string `json:"_id,omitempty"`
	Revf          string `json:"_rev,omitempty"`
	Keyf          string `json:"_key,omitempty"`
}

func (e *arangoError) IsError() bool {
	return e.IsErrorf
}

func (e *arangoError) Code() int {
	return e.Codef
}

func (e *arangoError) ErrorNum() int {
	return e.ErrorNumf
}

func (e *arangoError) ErrorMessage() string {
	return e.ErrorMessagef
}

func (e *arangoError) Error() string {
	return fmt.Sprintf("Code: %d, ErrorNum: %d, Message: %s\n",
		e.Code(),
		e.ErrorNum(),
		e.ErrorMessage())
}

func (e *arangoError) Id() string {
	return e.Idf
}
func (e *arangoError) Rev() string {
	return e.Revf
}
func (e *arangoError) Key() string {
	return e.Keyf
}

func newArangoError(code int, message string) ArangoError {
	return &arangoError{
		IsErrorf:      true,
		Codef:         code,
		ErrorNumf:     code,
		ErrorMessagef: message,
	}
}
