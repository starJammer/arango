package arango

import (
	gr "github.com/starJammer/grestclient"
	"net/http"
)

type database struct {
	connection Connection
	client     gr.Client
	name       string
}

func (d *database) Name() string {
	return d.name
}

func (d *database) CollectionEndpoint() CollectionEndpoint {
	cl := &collectionEndpoint{}
	cl.client = d.client.Clone()
	cl.client.BaseUrl().Path += CollectionPath
	cl.database = d

	return cl
}

func (d *database) DocumentEndpoint() DocumentEndpoint {
	doc := &documentEndpoint{}
	doc.client = d.client.Clone()
	doc.client.BaseUrl().Path += DocumentPath

	return doc
}

func (d *database) EdgeEndpoint() EdgeEndpoint {
	edge := &edgeEndpoint{}
	edge.client = d.client.Clone()
	edge.client.BaseUrl().Path += EdgePath

	return edge
}

func (d *database) Connection() Connection {
	return d.connection
}

func (d *database) Get() ([]string, error) {
	var result struct {
		Result []string `json:"result"`
	}
	var errorResult = &arangoError{}

	h, err := d.client.Get(
		DatabasePath,
		nil,
		nil,
		gr.UnmarshalMap{
			http.StatusOK:         gr.UnmarshalList(&result),
			http.StatusBadRequest: gr.UnmarshalList(errorResult),
			http.StatusForbidden:  gr.UnmarshalList(errorResult),
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return result.Result, nil

}

func (d *database) GetUser() ([]string, error) {

	var result struct {
		Result []string `json:"result"`
	}

	var errorResult = &arangoError{}

	h, err := d.client.Get(
		DatabasePath+"/user",
		nil,
		nil,
		gr.UnmarshalMap{
			http.StatusOK:         gr.UnmarshalList(&result),
			http.StatusBadRequest: gr.UnmarshalList(errorResult),
			http.StatusNotFound:   gr.UnmarshalList(errorResult),
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return result.Result, nil
}

type databaseDescriptor struct {
	Namef     string `json:"name"`
	Idf       string `json:"id"`
	Pathf     string `json:"path"`
	IsSystemf bool   `json:"isSystem"`
}

func (cr *databaseDescriptor) Name() string {
	return cr.Namef
}

func (cr *databaseDescriptor) Id() string {
	return cr.Idf
}

func (cr *databaseDescriptor) Path() string {
	return cr.Pathf
}

func (cr *databaseDescriptor) IsSystem() bool {
	return cr.IsSystemf
}

func (d *database) GetCurrent() (DatabaseDescriptor, error) {

	var result struct {
		Result *databaseDescriptor `json:"result"`
	}
	var errorResult = &arangoError{}

	errList := gr.UnmarshalList(errorResult)
	h, err := d.client.Get(
		DatabasePath+"/current",
		nil,
		nil,
		gr.UnmarshalMap{
			http.StatusOK:         gr.UnmarshalList(&result),
			http.StatusBadRequest: errList,
			http.StatusNotFound:   errList,
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return result.Result, nil
}

func (d *database) Post(name string, opts *PostDatabaseOptions) error {

	var errorResult = &arangoError{}
	if opts == nil {
		opts = new(PostDatabaseOptions)
	}
	opts.Name = name

	errList := gr.UnmarshalList(errorResult)
	h, err := d.client.Post(
		DatabasePath,
		nil,
		nil,
		opts,
		gr.UnmarshalMap{
			http.StatusBadRequest: errList,
			http.StatusNotFound:   errList,
			http.StatusConflict:   errList,
		},
	)

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusCreated {
		return errorResult
	}

	return nil
}

func (d *database) Delete(name string) error {

	var errorResult = &arangoError{}

	errList := gr.UnmarshalList(errorResult)
	h, err := d.client.Delete(
		DatabasePath+"/"+name,
		nil,
		nil,
		gr.UnmarshalMap{
			http.StatusBadRequest: errList,
			http.StatusNotFound:   errList,
		},
	)

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusOK {
		return errorResult
	}

	return nil
}
