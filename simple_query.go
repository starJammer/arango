package arango

import (
	"fmt"
)

type ByExampleQuery struct {
	Collection string      `json:"collection"`
	Example    interface{} `json:"example"`
	Skip       int         `json:"skip,omitempty"`
	Limit      int         `json:"limit,omitempty"`
	BatchSize  int         `json:"batchSize,omitempty"`
}

type FirstExampleQuery struct {
	Collection string      `json:"collection"`
	Example    interface{} `json:"example"`
}

func (db *Database) ByExampleQuery(query *ByExampleQuery) (*Cursor, error) {

	var c = new(Cursor)
	var e ArangoError

	endpoint := fmt.Sprintf("%s/simple/by-example",
		db.serverUrl.String(),
	)

	response, err := db.session.Put(endpoint, query, &c.json, &e)

	if err != nil {
		return nil, newError(err.Error())
	}

	switch response.Status() {
	case 201:
		c.db = db
		return c, nil
	default:
		return nil, e
	}

}

type firstExampleResult struct{
    Document interface{} `json:"document"`
    Error bool `json:"error"`
    Code int `json:"code"`
}

//FirstExample will call the PUT /_api/simple/first-example endpoint.
//The value pointed to by document is populated with the result from Arango.
func (db *Database) FirstExample(query *FirstExampleQuery, document interface{}) error {
	var e ArangoError
	endpoint := fmt.Sprintf("%s/simple/first-example",
		db.serverUrl.String(),
	)

    var result = &firstExampleResult{
        Document : document,
    }

	response, err := db.session.Put(endpoint, query, result, &e)

	if err != nil {
		return newError(err.Error())
	}

	switch response.Status() {
	case 200:
		return nil
	default:
		return e
	}
}
