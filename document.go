package arango

import (
	"fmt"
	gr "github.com/starJammer/grestclient"
	"net/http"
	"net/url"
)

type documentEndpoint struct {
	client gr.Client
}

func (doc *documentEndpoint) GetDocuments(
	collection string,
	returnType string,
) ([]string, error) {

	if returnType == "" {
		returnType = "path"
	}
	var errorResult = &arangoError{}
	var result struct {
		Documents []string `json:"documents"`
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

	return result.Documents, nil
}

func (doc *documentEndpoint) PostDocument(
	document interface{},
	collection string,
	options *PostDocumentOptions,
) error {

	var errorResult = &arangoError{}

	var query = url.Values{}
	query.Add("collection", collection)
	if options != nil {
		query.Add("createCollection", fmt.Sprintf("%t", options.CreateCollection))
		query.Add("waitForSync", fmt.Sprintf("%t", options.WaitForSync))
	}

	h, err := doc.client.Post(
		"",
		nil,
		query,
		document,
		document,
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

func (doc *documentEndpoint) GetDocument(documentHandle string, documentReceiver interface{}, options *GetDocumentOptions) error {

	var headers http.Header
	var query url.Values
	if options != nil {
		if options.IfNoneMatch != "" {
			headers.Add("If-None-Match", options.IfNoneMatch)
		}

		if options.Rev != "" {
			query.Add("rev", options.Rev)
		}
		if options.IfMatch != "" {
			headers.Add("If-Match", options.IfMatch)
		}
	}

	var errorResult = &arangoError{}

	h, err := doc.client.Get(
		fmt.Sprintf("/%s", documentHandle),
		headers,
		nil,
		documentReceiver,
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
