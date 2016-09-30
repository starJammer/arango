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

	h, err := s.client.Put(&gr.Params{
		Path: "/by-example",
		Body: body,
		UnmarshalMap: gr.UnmarshalMap{
			http.StatusCreated:    &cursor,
			http.StatusBadRequest: &errorResult,
			http.StatusForbidden:  &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	})

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusCreated {
		return nil, errorResult
	}

	return &Cursor{cursor: cursor, ce: s.Database().CursorEndpoint()}, nil
}

//PutSimpleAllKeys -> GET on /_api/simple/all-keys
//If opts.returnType is "" then the default should be used.
//Default is "path"
func (doc *SimpleEndpoint) PutSimpleAllKeys(opts *PutSimpleAllKeysOptions) ([]string, error) {

	var returnType string

	if opts != nil {
		returnType = opts.ReturnType
	}

	if returnType == "" {
		returnType = "path"
	}

	var errorResult = ArangoError{}
	var result struct {
		Documents []string `json:"result"`
	}

	h, err := doc.client.Put(&gr.Params{
		Path: "/all-keys",
		Body: opts,
		UnmarshalMap: gr.UnmarshalMap{
			http.StatusOK:         &result,
			http.StatusCreated:    &result,
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	})

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return nil, nil
}

//FirstExample -> PUT on /_api/simple/first-example
//result is where the returned document will be unmarshalled into
func (s *SimpleEndpoint) FirstExample(collection string, example interface{}, result interface{}) error {
	var body byExampleObject
	body.Collection = collection
	body.Example = example

	var errorResult = ArangoError{}

	var response struct {
		Document interface{} `json:"document"`
	}

	response.Document = result

	h, err := s.client.Put(&gr.Params{
		Path: "/first-example",
		Body: body,
		UnmarshalMap: gr.UnmarshalMap{
			http.StatusOK:         &response,
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	})

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusOK {
		return errorResult
	}

	return nil
}
