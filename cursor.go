package arango

import (
    "encoding/json"
)

//Cursor represents a collection of items that you can iterate over.
//This library will produce a Cursor type whenever ArangoDb produces
//a list of objects that match a query.     
type Cursor struct {
    db *Database
    json cursorResult
}

type cursorResult struct {
    Result []json.RawMessage `json:"result"`
    HasMore bool `json:"hasMore"`
    Count int `json:"count"`
    Error bool `json:"error"`
    Code int `json:"code"`
    Id string `json:"id"`
}

func (c Cursor) HasMore() bool {
    return len(c.json.Result) > 0 || c.json.HasMore
}

func (c Cursor) Count() int {
    return c.json.Count
}

func (c Cursor) Error() bool {
    return c.json.Error
}

func (c Cursor) Code() int {
    return c.json.Code
}

func (c Cursor) Next( next interface{} ) error {

    if len(c.json.Result) > 0 {

    } else if c.json.Id != "" {
    
    }

    return newError( "You called Next on a cursor that is invalid or doesn't have anymore results to return." )
}
