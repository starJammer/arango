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

func (db *Database) ByExampleQuery(query *ByExampleQuery) (*Cursor, error) {

	var c = new(Cursor)
	var e ArangoError

	endpoint := fmt.Sprintf("%s/simple/by-example",
		db.serverUrl.String(),
	)

	response, err := db.session.Put(endpoint, query, c, &e)

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
