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

//used to receive the old Values
//when ReturnOld params are set to true
type oldReceiver struct {
	OldReceiver interface{} `json:"old"`
}

//used to receive the new Values
//when ReturnNew params ar set to true
type newReceiver struct {
	NewReceiver interface{} `json:"new"`
}

type oldNewReceiver struct {
	OldReceiver interface{} `json:"old"`
	NewReceiver interface{} `json:"new"`
}

type DeleteDocumentOptions struct {
	Handle      string
	WaitForSync bool

	IfMatch string

	OldReceiver interface{}

	returnOld bool
}

func (de *DocumentEndpoint) DeleteDocument(opts *DeleteDocumentOptions) error {
	var errorResult = ArangoError{}
	var headers http.Header
	var query = make(url.Values)

	if opts == nil {
		opts = &DeleteDocumentOptions{}
	}

	var unmarshalMap = gr.UnmarshalMap{
		http.StatusBadRequest:         &errorResult,
		http.StatusNotFound:           &errorResult,
		http.StatusConflict:           &errorResult,
		http.StatusPreconditionFailed: &errorResult,
	}

	if opts.IfMatch != "" {
		headers = make(http.Header)
		if opts.IfMatch != "" {
			headers.Add("If-Match", opts.IfMatch)
		}
	}

	if opts.OldReceiver != nil {
		opts.returnOld = true
	}

	query.Add("waitForSync", fmt.Sprintf("%t", opts.WaitForSync))
	query.Add("returnOld", fmt.Sprintf("%t", opts.returnOld))

	if opts.returnOld {
		var t = &oldReceiver{
			OldReceiver: opts.OldReceiver,
		}
		unmarshalMap[http.StatusOK] = t
		unmarshalMap[http.StatusAccepted] = t
	} else {
		unmarshalMap[http.StatusOK] = opts.OldReceiver
		unmarshalMap[http.StatusAccepted] = opts.OldReceiver
	}

	var params = &gr.Params{
		Path:         "/" + opts.Handle,
		Query:        query,
		Headers:      headers,
		UnmarshalMap: unmarshalMap,
	}

	h, err := de.client.Delete(params)

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusOK &&
		h.StatusCode != http.StatusAccepted {
		return errorResult
	}

	return nil
}

type DeleteMultiDocumentsOptions struct {
	Handles     []interface{}
	Collection  string
	WaitForSync bool
	ReturnOld   bool
	IgnoreRevs  bool
	OldReceiver []interface{}
}

func (de *DocumentEndpoint) DeleteMultiDocuments(opts *DeleteMultiDocumentsOptions) error {
	var errorResult = ArangoError{}
	var query = make(url.Values)

	if opts == nil {
		opts = &DeleteMultiDocumentsOptions{}
	}

	var unmarshalMap = gr.UnmarshalMap{
		http.StatusBadRequest:         &errorResult,
		http.StatusNotFound:           &errorResult,
		http.StatusConflict:           &errorResult,
		http.StatusPreconditionFailed: &errorResult,
	}

	query.Add("waitForSync", fmt.Sprintf("%t", opts.WaitForSync))
	query.Add("returnOld", fmt.Sprintf("%t", opts.ReturnOld))
	query.Add("ignoreRevs", fmt.Sprintf("%t", opts.IgnoreRevs))

	if opts.OldReceiver != nil {
		if opts.ReturnOld {
			var t = make([]interface{}, len(opts.OldReceiver))
			for i, v := range opts.OldReceiver {
				t[i] = &oldReceiver{
					OldReceiver: v,
				}
			}
			unmarshalMap[http.StatusOK] = &t
			unmarshalMap[http.StatusAccepted] = &t
		} else {
			unmarshalMap[http.StatusOK] = &opts.OldReceiver
			unmarshalMap[http.StatusAccepted] = &opts.OldReceiver
		}
	}

	var params = &gr.Params{
		Path:         "/" + opts.Collection,
		Body:         opts.Handles,
		Query:        query,
		UnmarshalMap: unmarshalMap,
	}

	h, err := de.client.Delete(params)

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusOK &&
		h.StatusCode != http.StatusAccepted {
		return errorResult
	}

	return nil
}

type PatchDocumentOptions struct {
	Handle   string
	Document interface{}

	KeepNull     bool
	MergeObjects bool
	WaitForSync  bool
	IfMatch      string
	IgnoreRevs   bool

	//OldReiver specifies that old should be return
	OldReceiver interface{}
	//NewReceiver specifies that new should be return
	NewReceiver interface{}

	returnOld bool
	returnNew bool
}

func DefaultPatchDocumentOptions() *PatchDocumentOptions {
	return &PatchDocumentOptions{
		KeepNull:     true,
		IgnoreRevs:   true,
		MergeObjects: true,
	}
}

func (de *DocumentEndpoint) PatchDocument(opts *PatchDocumentOptions) error {

	if opts == nil {
		opts = DefaultPatchDocumentOptions()
	}

	var errorResult = ArangoError{}
	var query = make(url.Values)
	var headers http.Header

	if opts.OldReceiver != nil {
		opts.returnOld = true
	} else {
		opts.returnOld = false
	}

	if opts.NewReceiver != nil {
		opts.returnNew = true
	} else {
		opts.returnNew = false
	}

	query.Add("waitForSync", fmt.Sprintf("%t", opts.WaitForSync))
	query.Add("ignoreRevs", fmt.Sprintf("%t", opts.IgnoreRevs))
	query.Add("mergeObjects", fmt.Sprintf("%t", opts.MergeObjects))
	query.Add("keepNull", fmt.Sprintf("%t", opts.KeepNull))
	query.Add("returnOld", fmt.Sprintf("%t", opts.returnOld))
	query.Add("returnNew", fmt.Sprintf("%t", opts.returnNew))

	if opts.IfMatch != "" {
		headers = make(http.Header)
		if opts.IfMatch != "" {
			headers.Add("If-Match", opts.IfMatch)
		}
	}

	var unmarshalMap = gr.UnmarshalMap{
		http.StatusBadRequest:         &errorResult,
		http.StatusNotFound:           &errorResult,
		http.StatusPreconditionFailed: &errorResult,
	}

	if opts.returnOld && opts.returnNew {
		var t = &oldNewReceiver{
			OldReceiver: opts.OldReceiver,
			NewReceiver: opts.NewReceiver,
		}
		unmarshalMap[http.StatusCreated] = t
		unmarshalMap[http.StatusAccepted] = t
	} else if opts.returnOld {
		var t = &oldReceiver{
			OldReceiver: opts.OldReceiver,
		}
		unmarshalMap[http.StatusCreated] = t
		unmarshalMap[http.StatusAccepted] = t
	} else if opts.returnNew {
		var t = &newReceiver{
			NewReceiver: opts.NewReceiver,
		}
		unmarshalMap[http.StatusCreated] = t
		unmarshalMap[http.StatusAccepted] = t
	}

	params := &gr.Params{
		Path:         "/" + opts.Handle,
		Body:         opts.Document,
		Query:        query,
		Headers:      headers,
		UnmarshalMap: unmarshalMap,
	}

	h, err := de.client.Patch(params)

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusCreated &&
		h.StatusCode != http.StatusAccepted {
		return errorResult
	}

	return nil
}
