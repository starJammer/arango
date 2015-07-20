package arango

import (
	gr "github.com/starJammer/grestclient"
	"net/url"
)

type documentEndpoint struct {
	client gr.Client
}

func (doc *documentEndpoint) GetDocuments(collection, returnType string) ([]string, error) {

	if returnType == "" {
		returnType = "path"
	}
	var errorResult = &arangoError{}
	var result struct {
		Documents []string `json:"documents"`
	}

	h, err := doc.client.Get("",
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

func (doc *documentEndpoint) PostDocument(document interface{}) error {

	return nil
}
