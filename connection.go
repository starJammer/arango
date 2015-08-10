package arango

import (
	"fmt"
	gr "github.com/starJammer/grestclient"
	"net/http"
	"net/url"
)

type connection struct {
	client gr.Client
}

//NewConnection creates a new Connection that will
//use the given url as the address for the arango
//server. The url should not contain anything
//other than basic authentication info and
//the server's base url. If you include a path it will
//be stripped/ignored.
//Ex. http://localhost:8529.
func NewConnection(serverUrl *url.URL) (Connection, error) {

	c := &connection{}

	client, err := gr.New(serverUrl)
	if err != nil {
		return nil, err
	}
	gr.SetupForJson(client)

	c.client = client

	return c, nil
}

type version struct {
	S string            `json:"server"`
	V string            `json:"version"`
	D map[string]string `json:"details"`
}

func (v *version) Version() string {
	return v.V
}

func (v *version) Server() string {
	return v.S
}

func (v *version) Details() map[string]string {
	return v.D
}

func (c *connection) Version(details bool) (Version, error) {
	v := &version{}

	params := url.Values{}
	if details {
		params.Add("details", "true")
	}
	h, err := c.client.Get(
		VersionPath,
		nil,
		params,
		gr.UnmarshalMap{
			http.StatusOK: gr.UnmarshalList(v),
		},
	)

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, newArangoError(h.StatusCode, "Unxpected error fetching version.")
	}

	return v, nil
}

func (c *connection) Database(name string) Database {
	db := &database{}
	db.connection = c
	db.name = name
	db.client = c.client.Clone()
	db.client.BaseUrl().Path += fmt.Sprintf(DatabasePrefix, name)

	return db
}

func (c *connection) GetGrestClient() gr.Client {
	return c.client
}
