package arango

import (
	"encoding/json"
	"fmt"
)

//Cursor represents a collection of items that you can iterate over.
//This library will produce a Cursor type whenever ArangoDb produces
//a list of objects that match a query.
type Cursor struct {
	db   *Database
	json cursorResult
}

type cursorResult struct {
	Result  []json.RawMessage `json:"result"`
	HasMore bool              `json:"hasMore"`
	Count   int               `json:"count"`
	Error   bool              `json:"error"`
	Code    int               `json:"code"`
	Id      string            `json:"id"`
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

//Next retrieves the next item from the cursor.
//According to the arango docs :
//Note that even if hasMore returns true, the next 
//call might still return no documents. 
//If, however, hasMore is false, then the cursor 
//is exhausted. Once the hasMore attribute has a value 
//of false, the client can stop. 
func (c *Cursor) Next(next interface{}) error {

	if len(c.json.Result) > 0 {
		err := json.Unmarshal(c.json.Result[0], next)
		if err != nil {
			return newError(err.Error())
		}
		c.json.Result = c.json.Result[1:len(c.json.Result)]
		return nil
	} else if c.json.Id != "" {
		endpoint := fmt.Sprintf("%s/cursor/%s",
			c.db.serverUrl.String(),
			c.json.Id,
		)

		var e ArangoError
		response, err := c.db.session.Put(endpoint, nil, &c.json, &e)

		if err != nil {
			return newError(err.Error())
		}
		switch response.Status() {
		case 200:
			if len(c.json.Result) > 0 {
				err := json.Unmarshal(c.json.Result[0], next)
				if err != nil {
					return newError(err.Error())
				}
				c.json.Result = c.json.Result[1:len(c.json.Result)]
			}
			return nil
		default:
			return e
		}
	}
	return newError("You called Next on a cursor that is invalid or doesn't have anymore results to return.")
}

func (c Cursor) Close() error {
	endpoint := fmt.Sprintf("%s/cursor/%s",
		c.db.serverUrl.String(),
		c.json.Id,
	)
	var e ArangoError
	response, err := c.db.session.Delete(endpoint, nil, &e)

	if err != nil {
		return newError(err.Error())
	}

	switch response.Status() {
	case 202:
		return nil
	default:
		return e
	}
}
