package arango

import (
	"fmt"
	gr "github.com/starJammer/grestclient"
	"net/http"
	"net/url"
	"strings"
)

type EdgeEndpoint struct {
	client   *gr.Client
	database *Database
}

type GetEdgesOptions GetDocumentsOptions

func (e *EdgeEndpoint) Database() *Database {
	return e.database
}

//GetEdges -> GET on /_api/edge
//If pased in returnType is "" then the default should be used.
//Default is "path"
func (e *EdgeEndpoint) GetEdges(
	collection string,
	opts *GetEdgesOptions,
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
		Edges []string `json:"documents"`
	}

	h, err := e.client.Get(
		"",
		nil,
		url.Values{
			"collection": []string{collection},
			"type":       []string{returnType},
		},
		gr.UnmarshalMap{
			http.StatusOK:         &result,
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return result.Edges, nil
}

type PostEdgeOptions PostDocumentOptions

//PostEdge -> POST on /_api/edge
//PostEdgeOptions are optional.
//The same edge is populated with _id, _key, _rev attributes
//if possible on a successful POST of the edge.
func (e *EdgeEndpoint) PostEdge(
	edge interface{},
	collection string,
	from string,
	to string,
	options *PostEdgeOptions,
) error {

	var errorResult = ArangoError{}

	var query = url.Values{}
	query.Add("collection", collection)
	query.Add("from", from)
	query.Add("to", to)
	if options != nil {
		query.Add("createCollection", fmt.Sprintf("%t", options.CreateCollection))
		query.Add("waitForSync", fmt.Sprintf("%t", options.WaitForSync))
	}

	h, err := e.client.Post(
		"",
		nil,
		query,
		edge,
		gr.UnmarshalMap{
			http.StatusCreated:    edge,
			http.StatusAccepted:   edge,
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	)

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusCreated &&
		h.StatusCode != http.StatusAccepted {
		return errorResult
	}

	if e, ok := edge.(Edge); ok {
		e.SetFrom(from)
		e.SetTo(to)
	}

	return nil
}

type GetEdgeOptions GetDocumentOptions

//GetEdge -> GET on /_api/edge/{edge-handle}
//edgeReceiver is where the edge will be json.Unmarshaled into.
//GetEdgeOptions are optional and can be nil.
//EdgeReceiver cannot be populated if you provide an
//If-None-Match option and the server returns a 304 because the edge
//revision matches If-None-Match. See arango docs for more info.
//In this case, no error is returned and the edgeReceiver isn't
//altered at all because the server didn't return any edge
//attributes.
//
//EdgeReceiver IS NOT populated if you provide an
//If-Match option and the server returns a 412 because the edge
//revision does not match If-Match.  In this case, use the
//error that is returned to get the latest revision number
//by calling error.Rev()
func (e *EdgeEndpoint) GetEdge(edgeHandle string, edgeReceiver interface{}, options *GetEdgeOptions) error {

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

	h, err := e.client.Get(
		fmt.Sprintf("/%s", edgeHandle),
		headers,
		query,
		gr.UnmarshalMap{
			http.StatusOK:                 edgeReceiver,
			http.StatusBadRequest:         &errorResult,
			http.StatusNotFound:           &errorResult,
			http.StatusPreconditionFailed: &errorResult,
		},
	)

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusOK && h.StatusCode != http.StatusNotModified {
		return errorResult
	}

	return nil
}

type HeadEdgeOptions HeadDocumentOptions

//HeadEdge -> HEAD on /_api/edge/{edge-handle}
//Returns the current revision of the edge.
//If the edge doesn't exist the returned revision is blank.
//In all other cases, the current revision is returned and
//the error is nil.
func (e *EdgeEndpoint) HeadEdge(edgeHandle string, options *HeadEdgeOptions) (string, error) {

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

	h, err := e.client.Head(
		fmt.Sprintf("/%s", edgeHandle),
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

type PutEdgeOptions PutDocumentOptions

//PutEdge -> PUT on /_api/edge/{edge-handle}
func (e *EdgeEndpoint) PutEdge(edgeHandle string, edge interface{}, options *PutEdgeOptions) error {

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

	h, err := e.client.Put(
		fmt.Sprintf("/%s", edgeHandle),
		headers,
		query,
		edge, //edge is the body
		gr.UnmarshalMap{
			http.StatusCreated:            edge,
			http.StatusAccepted:           edge,
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

type PatchEdgeOptions PatchDocumentOptions

func DefaultPatchEdgeOptions() *PatchEdgeOptions {
	return &PatchEdgeOptions{
		KeepNull:     true,
		MergeObjects: true,
		WaitForSync:  false,
	}
}

//PatchEdge -> PATCH on /_api/edge/{edge-handle}
func (e *EdgeEndpoint) PatchEdge(edgeHandle string, edge interface{}, options *PatchEdgeOptions) error {

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

	h, err := e.client.Patch(
		fmt.Sprintf("/%s", edgeHandle),
		headers,
		query,
		//edge is the body
		edge,
		gr.UnmarshalMap{
			http.StatusCreated:            edge,
			http.StatusAccepted:           edge,
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

type DeleteEdgeOptions DeleteDocumentOptions

//DeleteEdge -> DELETE on /_api/edge/{edge-handle}
func (e *EdgeEndpoint) DeleteEdge(edgeHandle string, options *DeleteEdgeOptions) error {

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

	h, err := e.client.Delete(
		fmt.Sprintf("/%s", edgeHandle),
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

type Edge interface {
	SetFrom(string)
	SetTo(string)
}

func (e *EdgeImplementation) From() string {
	return e.ArangoFrom
}

func (e *EdgeImplementation) SetFrom(from string) {
	e.ArangoFrom = from
}

func (e *EdgeImplementation) To() string {
	return e.ArangoTo
}

func (e *EdgeImplementation) SetTo(to string) {
	e.ArangoTo = to
}
