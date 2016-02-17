package arango

import (
	"fmt"
	gr "github.com/starJammer/grestclient"
	"net/http"
	"net/url"
	"strings"
)

type DocumentEndpoint struct {
	client   *gr.Client
	database *Database
}

type GetDocumentsOptions struct {
	ReturnType string
}

func (doc *DocumentEndpoint) Database() *Database {
	return doc.database
}

//GetDocuments -> GET on /_api/document
//If pased in returnType is "" then the default should be used.
//Default is "path"
func (doc *DocumentEndpoint) GetDocuments(
	collection string,
	opts *GetDocumentsOptions,
) ([]string, error) {

	var returnType string

	if opts != nil {
		returnType = opts.ReturnType
	}

	if returnType == "" {
		returnType = "path"
	}

	var errorResult = ArangoError{}
	var result struct {
		Documents []string `json:"documents"`
	}

	h, err := doc.client.Get(&gr.Request{
		Path: "",
		Query: url.Values{
			"collection": []string{collection},
			"type":       []string{returnType},
		},
		UnmarshalMap: gr.UnmarshalMap{
			http.StatusOK:         &result,
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

	return result.Documents, nil
}

type PostDocumentOptions struct {
	CreateCollection bool //default is false
	WaitForSync      bool //default is false
}

//PostDocument -> POST on /_api/document
//PostDocumentOptions are optional.
//The same document is populated with _id, _key, _rev attributes
//if possible on a successful POST of the document.
func (doc *DocumentEndpoint) PostDocument(
	document interface{},
	collection string,
	options *PostDocumentOptions,
) error {

	var errorResult = ArangoError{}

	var query = url.Values{}
	query.Add("collection", collection)
	if options != nil {
		query.Add("createCollection", fmt.Sprintf("%t", options.CreateCollection))
		query.Add("waitForSync", fmt.Sprintf("%t", options.WaitForSync))
	}

	h, err := doc.client.Post(&gr.Request{
		Path:    "",
		Headers: nil,
		Query:   query,
		Body:    document,
		UnmarshalMap: gr.UnmarshalMap{
			http.StatusCreated:    document,
			http.StatusAccepted:   document,
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	})

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusCreated &&
		h.StatusCode != http.StatusAccepted {
		return errorResult
	}

	return nil
}

type GetDocumentOptions struct {
	//Rev is used in the query
	Rev string
	//IfMatch is a header equivalent of IfMatch
	IfMatch string

	//IfNoneMatch is header
	IfNoneMatch string
}

//GetDocument -> GET on /_api/document/{document-handle}
//documentReceiver is where the document will be json.Unmarshaled into.
//GetDocumentOptions are optional and can be nil.
//DocumentReceiver cannot be populated if you provide an
//If-None-Match option and the server returns a 304 because the document
//revision matches If-None-Match. See arango docs for more info.
//In this case, no error is returned and the documentReceiver isn't
//altered at all because the server didn't return any document
//attributes.
//
//DocumentReceiver IS NOT populated if you provide an
//If-Match option and the server returns a 412 because the document
//revision does not match If-Match.  In this case, use the
//error that is returned to get the latest revision number
//by calling error.Rev()
func (doc *DocumentEndpoint) GetDocument(
	documentHandle string,
	documentReceiver interface{},
	options *GetDocumentOptions,
) error {

	var headers http.Header
	var query url.Values
	if options != nil {
		if options.Rev != "" {
			query = make(url.Values)
			query.Add("rev", options.Rev)
		}

		headers = make(http.Header)
		if options.IfNoneMatch != "" {
			headers.Add("If-None-Match", options.IfNoneMatch)
		}

		if options.IfMatch != "" {
			headers.Add("If-Match", options.IfMatch)
		}
	}

	var errorResult = ArangoError{}

	h, err := doc.client.Get(&gr.Request{
		Path:    fmt.Sprintf("/%s", documentHandle),
		Headers: headers,
		Query:   query,
		UnmarshalMap: gr.UnmarshalMap{
			http.StatusOK:                 documentReceiver,
			http.StatusBadRequest:         &errorResult,
			http.StatusNotFound:           &errorResult,
			http.StatusPreconditionFailed: &errorResult,
		},
	})

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusOK && h.StatusCode != http.StatusNotModified {
		return errorResult
	}

	return nil
}

type HeadDocumentOptions GetDocumentOptions

//HeadDocument -> HEAD on /_api/document/{document-handle}
//Returns the current revision of the document.
//If the document doesn't exist the returned revision is blank.
//In all other cases, the current revision is returned and
//the error is nil.
func (doc *DocumentEndpoint) HeadDocument(documentHandle string, options *HeadDocumentOptions) (string, error) {

	var headers http.Header
	var query url.Values
	if options != nil {
		if options.Rev != "" {
			query = make(url.Values)
			query.Add("rev", options.Rev)
		}

		headers = make(http.Header)
		if options.IfNoneMatch != "" {
			headers.Add("If-None-Match", options.IfNoneMatch)
		}

		if options.IfMatch != "" {
			headers.Add("If-Match", options.IfMatch)
		}
	}

	h, err := doc.client.Head(
		fmt.Sprintf("/%s", documentHandle),
		headers,
		query,
	)

	if err != nil {
		return "", err
	}

	if h.StatusCode == http.StatusBadRequest {
		return "", newArangoError(h.StatusCode, "Malformed request.")
	}

	var revision = strings.Trim(h.Header.Get("Etag"), "\"")
	if h.StatusCode == http.StatusNotModified {
		return revision, newArangoError(h.StatusCode, "Document hasn't been modified.")
	}

	if h.StatusCode == http.StatusPreconditionFailed {
		return revision, newArangoError(h.StatusCode, "Document's revision different from If-None-Match.")
	}

	if h.StatusCode != http.StatusOK {
		return "", newArangoError(h.StatusCode, "Unknown response from arango.")
	}

	return revision, nil
}

type PutDocumentOptions struct {
	WaitForSync bool
	Rev         string
	Policy      Policy
	IfMatch     string
}

//PutDocument -> PUT on /_api/document/{document-handle}
func (doc *DocumentEndpoint) PutDocument(documentHandle string, document interface{}, options *PutDocumentOptions) error {

	var headers http.Header
	var query url.Values
	if options != nil {
		query = make(url.Values)
		query.Add("waitForSync", fmt.Sprintf("%t", options.WaitForSync))

		if options.Rev != "" {
			query.Add("rev", options.Rev)
		}
		if options.Policy != "" {
			query.Add("policy", string(options.Policy))
		}

		headers = make(http.Header)

		if options.IfMatch != "" {
			headers.Add("If-Match", options.IfMatch)
		}

	}

	var errorResult = ArangoError{}

	h, err := doc.client.Put(
		fmt.Sprintf("/%s", documentHandle),
		headers,
		query,
		document, //document is the body
		gr.UnmarshalMap{
			http.StatusCreated:            document,
			http.StatusAccepted:           document,
			http.StatusBadRequest:         &errorResult,
			http.StatusNotFound:           &errorResult,
			http.StatusPreconditionFailed: &errorResult,
		},
	)

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusCreated &&
		h.StatusCode != http.StatusAccepted {
		return errorResult
	}

	return nil

}

type PatchDocumentOptions struct {
	KeepNull     bool
	MergeObjects bool

	WaitForSync bool
	Rev         string
	Policy      Policy
	IfMatch     string
}

func DefaultPatchDocumentOptions() *PatchDocumentOptions {
	return &PatchDocumentOptions{
		KeepNull:     true,
		MergeObjects: true,
		WaitForSync:  false,
	}
}

//PatchDocument -> PATCH on /_api/document/{document-handle}
func (doc *DocumentEndpoint) PatchDocument(documentHandle string, document interface{}, options *PatchDocumentOptions) error {

	var headers http.Header
	var query url.Values
	if options != nil {
		query = make(url.Values)
		query.Add("keepNull", fmt.Sprintf("%t", options.KeepNull))
		query.Add("mergeObjects", fmt.Sprintf("%t", options.MergeObjects))
		query.Add("waitForSync", fmt.Sprintf("%t", options.WaitForSync))

		if options.Rev != "" {
			query.Add("rev", options.Rev)
		}
		if options.Policy != "" {
			query.Add("policy", string(options.Policy))
		}

		headers = make(http.Header)

		if options.IfMatch != "" {
			headers.Add("If-Match", options.IfMatch)
		}

	}

	var errorResult = ArangoError{}

	h, err := doc.client.Patch(
		fmt.Sprintf("/%s", documentHandle),
		headers,
		query,
		//document is the body
		document,
		gr.UnmarshalMap{
			//document is used as the successResult so it gets
			//populated with the new revision info
			http.StatusCreated:            document,
			http.StatusAccepted:           document,
			http.StatusBadRequest:         &errorResult,
			http.StatusNotFound:           &errorResult,
			http.StatusPreconditionFailed: &errorResult,
		},
	)

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusCreated &&
		h.StatusCode != http.StatusAccepted {
		return errorResult
	}

	return nil

}

type DeleteDocumentOptions struct {
	Rev         string
	Policy      Policy
	WaitForSync string
	IfMatch     string
}

//DeleteDocument -> DELETE on /_api/document/{document-handle}
func (doc *DocumentEndpoint) DeleteDocument(documentHandle string, options *DeleteDocumentOptions) error {

	var headers http.Header
	var query url.Values
	if options != nil {
		query = make(url.Values)
		query.Add("waitForSync", fmt.Sprintf("%t", options.WaitForSync))

		if options.Rev != "" {
			query.Add("rev", options.Rev)
		}
		if options.Policy != "" {
			query.Add("policy", string(options.Policy))
		}

		headers = make(http.Header)

		if options.IfMatch != "" {
			headers.Add("If-Match", options.IfMatch)
		}

	}

	var errorResult = ArangoError{}

	h, err := doc.client.Delete(
		fmt.Sprintf("/%s", documentHandle),
		headers,
		query,
		gr.UnmarshalMap{
			http.StatusBadRequest:         &errorResult,
			http.StatusNotFound:           &errorResult,
			http.StatusPreconditionFailed: &errorResult,
		},
	)

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusCreated &&
		h.StatusCode != http.StatusAccepted {
		return errorResult
	}

	return nil
}

//DocumentImplementation is an embeddable type that
//you can use to easily gain access to arango specific
//data. It will help capture the _id, _key, and _rev
//attributes from responses made by  arangodb
type DocumentImplementation struct {
	ArangoId  string `json:"_id,omitempty"`
	ArangoRev string `json:"_rev,omitempty"`
	ArangoKey string `json:"_key,omitempty"`
}

func (d *DocumentImplementation) Id() string {
	return d.ArangoId
}

func (d *DocumentImplementation) SetId(id string) {
	d.ArangoId = id
}

func (d *DocumentImplementation) Rev() string {
	return d.ArangoRev
}

func (d DocumentImplementation) SetRev(rev string) {
	d.ArangoRev = rev
}

func (d DocumentImplementation) Key() string {
	return d.ArangoKey
}

func (d DocumentImplementation) SetKey(key string) {
	d.ArangoKey = key
}
