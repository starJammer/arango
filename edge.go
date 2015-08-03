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
		&result, errorResult)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != 200 {
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
		edge,
		errorResult,
	)

	if err != nil {
		return err
	}

	if h.StatusCode != 201 && h.StatusCode != 202 {
		return errorResult
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
		edgeReceiver,
		errorResult,
	)

	if err != nil {
		return err
	}

	if h.StatusCode != 200 && h.StatusCode != 304 {
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

	if h.StatusCode == 400 {
		return "", newArangoError(h.StatusCode, "Malformed request.")
	}

	if h.StatusCode != 200 && h.StatusCode != 304 && h.StatusCode != 412 {
		return "", newArangoError(h.StatusCode, "Unknown response from arango.")
	}

	return strings.Trim(h.Header.Get("Etag"), "\""), nil
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
		edge, //edge is used as the successResult so it gets
		//populated with the new revision info
		errorResult,
	)

	if err != nil {
		return err
	}

	if h.StatusCode != 201 && h.StatusCode != 202 {
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
		//edge is used as the successResult so it gets
		//populated with the new revision info
		edge,
		errorResult,
	)

	if err != nil {
		return err
	}

	if h.StatusCode != 201 && h.StatusCode != 202 {
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
		nil,
		errorResult,
	)

	if err != nil {
		return err
	}

	if h.StatusCode != 201 && h.StatusCode != 202 {
		return errorResult
	}

	return nil
}
