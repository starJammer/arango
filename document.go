package arango

import (
	"fmt"
	gr "github.com/starJammer/grestclient"
	"net/http"
	"net/url"
)

//DocumentImplementation is an embeddable type that
//you can use to easily gain access to arango specific
//data. It will help capture the _id, _key, and _rev
//attributes from responses made by  arangodb
type DocumentImplementation struct {
	ArangoId  string `json:"_id,omitempty"`
	ArangoRev string `json:"_rev,omitempty"`
	ArangoKey string `json:"_key,omitempty"`
}

//EdgeImplementation is an embeddable type that
//you can use to easily gain access to arango
//specific attributes for edges. These include
//the _id, _key, and _rev attributes from the
//DocumentImplementation as well as the _to and
//_from attributes for edges.
type EdgeImplementation struct {
	DocumentImplementation
	ArangoFrom string `json:"_from,omitempty"`
	ArangoTo   string `json:"_to,omitempty"`
}

type DocumentEndpoint struct {
	client   *gr.Client
	database *Database
}

func (d *DocumentEndpoint) Database() *Database {
	return d.database
}

type GetDocumentOptions struct {
	Handle      string
	Document    interface{}
	IfNoneMatch string
	IfMatch     string
}

func (de *DocumentEndpoint) GetDocument(opts *GetDocumentOptions) error {

	var errorResult = ArangoError{}
	var headers http.Header

	if opts == nil {
		opts = &GetDocumentOptions{}
	}

	var unmarshalMap = gr.UnmarshalMap{
		http.StatusOK:         opts.Document,
		http.StatusBadRequest: &errorResult,
		http.StatusNotFound:   &errorResult,
		http.StatusConflict:   &errorResult,
	}

	if opts.IfNoneMatch != "" || opts.IfMatch != "" {
		headers = make(http.Header)
		if opts.IfNoneMatch != "" {
			headers.Add("If-None-Match", opts.IfNoneMatch)
		}
		if opts.IfMatch != "" {
			headers.Add("If-Match", opts.IfMatch)
		}
	}

	var params = &gr.Params{
		Path:         "/" + opts.Handle,
		Headers:      headers,
		UnmarshalMap: unmarshalMap,
	}

	h, err := de.client.Get(params)

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusOK {
		return errorResult
	}

	return nil
}

type PostDocumentOptions struct {
	//Use this if you want to post only one document
	Document interface{}
	//Use this if you want to post multiple documents
	//This takes precedence over the single document
	MultiDocuments []interface{}
	Collection     string
	WaitForSync    bool

	//ReturnNew doesn't currently do anything because
	//the same object you pass is used in the json.Unmarshal
	//call in the end anyway, so it would already have
	//all the data except for the arango meta, which is
	//populated by the json.Unmarshal call.
	ReturnNew bool
}

func (de *DocumentEndpoint) PostDocuments(opts *PostDocumentOptions) error {

	var errorResult = ArangoError{}
	var query = make(url.Values)

	if opts == nil {
		opts = &PostDocumentOptions{}
	}

	var unmarshalMap = gr.UnmarshalMap{
		http.StatusBadRequest: &errorResult,
		http.StatusNotFound:   &errorResult,
		http.StatusConflict:   &errorResult,
	}

	var params = &gr.Params{
		Path:         "/" + opts.Collection,
		UnmarshalMap: unmarshalMap,
	}

	query.Add("waitForSync", fmt.Sprintf("%t", opts.WaitForSync))
	query.Add("returnNew", fmt.Sprintf("%t", opts.ReturnNew))

	if opts.MultiDocuments != nil {
		params.Body = opts.MultiDocuments
		unmarshalMap[http.StatusCreated] = &opts.MultiDocuments
		unmarshalMap[http.StatusAccepted] = &opts.MultiDocuments
	} else if opts.Document != nil {
		params.Body = opts.Document
		unmarshalMap[http.StatusCreated] = opts.Document
		unmarshalMap[http.StatusAccepted] = opts.Document
	}

	h, err := de.client.Post(params)

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusCreated && h.StatusCode != http.StatusAccepted {
		return errorResult
	}

	return nil

}

func DeleteDocument() {

}

func HeadDocument() {

}

func PatchDocument() {

}
