package arango

import (
	"fmt"
	gr "github.com/starJammer/grestclient"
	"net/http"
	"net/url"
)

//Connection represents a RESTFUL gateway to an arangodb server
//Use NewConnection  to create it.
type Connection struct {
	client *gr.Client
}

//NewConnection creates a new Connection that will
//use the given url as the address for the arango
//server. The url should not contain anything
//other than basic authentication info and
//the server's base url. If you include a path it will
//be stripped/ignored.
//Ex. http://localhost:8529.
func NewConnection(serverUrl *url.URL) (*Connection, error) {

	c := &Connection{}

	client, err := gr.New(serverUrl)
	if err != nil {
		return nil, err
	}
	gr.SetupForJson(client)

	c.client = client

	return c, nil
}

type Version struct {
	Server  string            `json:"server"`
	Version string            `json:"version"`
	Details map[string]string `json:"details"`
}

func (c *Connection) Version(details bool) (*Version, error) {
	v := &Version{}

	params := url.Values{}
	if details {
		params.Add("details", "true")
	}
	h, err := c.client.Get(&gr.Request{
		Path:  VersionPath,
		Query: params,
		UnmarshalMap: gr.UnmarshalMap{
			http.StatusOK: v,
		},
	})

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, newArangoError(h.StatusCode, "Unxpected error fetching version.")
	}

	return v, nil
}

//Database returns a RESTFUL gateway to the database
//endpoint. The url used would be the url for the connection
//plus the adequate path for this database.
//Ex. http://localhost:8529/_db/{name}/_api/database where
//{name} is the passed in database name.
func (c *Connection) Database(name string) *Database {
	db := &Database{}
	db.connection = c
	db.name = name
	db.client = c.client.Clone()
	db.client.BaseUrl().Path += fmt.Sprintf(DatabasePrefix, name)

	return db
}

func (c *Connection) GetGrestClient() *gr.Client {
	return c.client
}
