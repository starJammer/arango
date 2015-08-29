package arango

import (
	"net/http"

	gr "github.com/starJammer/grestclient"
)

type SimpleEndpoint struct {
	client   *gr.Client
	database *Database
}

//Database gets the related database endpoint
//for this collection endpoint
func (s *SimpleEndpoint) Database() *Database {
	return s.database
}

type ByExampleOptions struct {
	Skip      int `json:"skip,omitempty"`
	Limit     int `json:"limit,omitempty"`
	BatchSize int `json:"batchSize,omitempty"`
}

type byExampleObject struct {
	Collection string      `json:"collection"`
	Example    interface{} `json:"example"`
	ByExampleOptions
}

func (s *SimpleEndpoint) ByExample(collection string, example interface{}, opts *ByExampleOptions) (*Cursor, error) {

	var body byExampleObject
	body.Collection = collection
	body.Example = example

	if opts != nil {
		body.ByExampleOptions = *opts
	}

	var errorResult = ArangoError{}

	var cursor cursor

	h, err := s.client.Put(
		"/by-example",
		nil,
		nil,
		body,
		gr.UnmarshalMap{
			http.StatusCreated:    &cursor,
			http.StatusBadRequest: &errorResult,
			http.StatusForbidden:  &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusCreated {
		return nil, errorResult
	}

	return &Cursor{cursor: cursor, ce: s.Database().CursorEndpoint()}, nil
}