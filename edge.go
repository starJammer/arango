package arango

import (
	"fmt"
	gr "github.com/starJammer/grestclient"
	"net/http"
	"net/url"
	"strings"
)

type edgeEndpoint struct {
	client gr.Client
}

func (doc *edgeEndpoint) GetEdges(
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

	var errorResult = &arangoError{}
	var result struct {
		Edges []string `json:"documents"`
	}

	h, err := doc.client.Get(
		"",
		nil,
		url.Values{
			"collection": []string{collection},
			"type":       []string{returnType},
		},
		gr.UnmarshalMap{
			http.StatusOK:         &result,
			http.StatusBadRequest: errorResult,
			http.StatusNotFound:   errorResult,
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

func (doc *edgeEndpoint) PostEdge(
	edge interface{},
	collection string,
	from string,
	to string,
	options *PostEdgeOptions,
) error {

	var errorResult = &arangoError{}

	var query = url.Values{}
	query.Add("collection", collection)
	query.Add("from", from)
	query.Add("to", to)
	if options != nil {
		query.Add("createCollection", fmt.Sprintf("%t", options.CreateCollection))
		query.Add("waitForSync", fmt.Sprintf("%t", options.WaitForSync))
	}

	h, err := doc.client.Post(
		"",
		nil,
		query,
		edge,
		gr.UnmarshalMap{
			http.StatusCreated:    edge,
			http.StatusAccepted:   edge,
			http.StatusBadRequest: errorResult,
			http.StatusNotFound:   errorResult,
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

func (doc *edgeEndpoint) GetEdge(edgeHandle string, edgeReceiver interface{}, options *GetEdgeOptions) error {

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

	var errorResult = &arangoError{}

	h, err := doc.client.Get(
		fmt.Sprintf("/%s", edgeHandle),
		headers,
		query,
		gr.UnmarshalMap{
			http.StatusOK:                 edgeReceiver,
			http.StatusBadRequest:         errorResult,
			http.StatusNotFound:           errorResult,
			http.StatusPreconditionFailed: errorResult,
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

func (doc *edgeEndpoint) HeadEdge(edgeHandle string, options *HeadEdgeOptions) (string, error) {

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

func (doc *edgeEndpoint) PutEdge(edgeHandle string, edge interface{}, options *PutEdgeOptions) error {

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

	var errorResult = &arangoError{}

	h, err := doc.client.Put(
		fmt.Sprintf("/%s", edgeHandle),
		headers,
		query,
		edge, //edge is the body
		gr.UnmarshalMap{
			http.StatusCreated:            edge,
			http.StatusAccepted:           edge,
			http.StatusBadRequest:         errorResult,
			http.StatusNotFound:           errorResult,
			http.StatusPreconditionFailed: errorResult,
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

func (doc *edgeEndpoint) PatchEdge(edgeHandle string, edge interface{}, options *PatchEdgeOptions) error {

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

	var errorResult = &arangoError{}

	h, err := doc.client.Patch(
		fmt.Sprintf("/%s", edgeHandle),
		headers,
		query,
		//edge is the body
		edge,
		gr.UnmarshalMap{
			http.StatusCreated:            edge,
			http.StatusAccepted:           edge,
			http.StatusBadRequest:         errorResult,
			http.StatusNotFound:           errorResult,
			http.StatusPreconditionFailed: errorResult,
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

func (doc *edgeEndpoint) DeleteEdge(edgeHandle string, options *DeleteEdgeOptions) error {

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

	var errorResult = &arangoError{}

	h, err := doc.client.Delete(
		fmt.Sprintf("/%s", edgeHandle),
		headers,
		query,
		gr.UnmarshalMap{
			http.StatusBadRequest:         errorResult,
			http.StatusNotFound:           errorResult,
			http.StatusPreconditionFailed: errorResult,
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
