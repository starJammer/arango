package arango

import (
	"fmt"
	gr "github.com/starJammer/grestclient"
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

type PostDocumentOptions struct {
	//Use this if you want to post only one document
	Document interface{}
	//Use this if you want to post multiple documents
	//This takes precedence over the single document
	Documents   []interface{}
	Collection  string
	WaitForSync bool
	ReturnNew   bool
}

func (de *DocumentEndpoint) GetDocument(opts *GetDocumentOptions) {

}

func (de *DocumentEndpoint) PostDocuments(opts *PostDocumentOptions) error {

	var errorResult = ArangoError{}
	var query url.Values

	if opts.Collection != "" {
		query.Add("collection", opts.Collection)
		query.Add("waitForSync", fmt.Sprintf("%t", opts.WaitForSync))
		query.Add("returnNew", fmt.Sprintf("%t", opts.ReturnNew))
	}

	h, err := de.client.Get(&gr.Params{
		Path:         "",
		Query:        query,
		UnmarshalMap: gr.UnmarshalMap{},
	})

	if err != nil {
		return nil, err
	}

}

func DeleteDocument() {

}

func HeadDocument() {

}

func PatchDocument() {

}
