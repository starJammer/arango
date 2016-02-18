package arango

import (
	gr "github.com/starJammer/grestclient"
	"net/http"
)

type Database struct {
	connection *Connection
	client     *gr.Client
	name       string

	ce   *CollectionEndpoint
	de   *DocumentEndpoint
	se   *SimpleEndpoint
	ee   *EdgeEndpoint
	curE *CursorEndpoint
}

//Name returns the name of the database that this endpoint accesses.
//In other words, what was Connection.Database called with?
func (d *Database) Name() string {
	return d.name
}

//Collection gets the collection endpoint.
func (d *Database) CollectionEndpoint() *CollectionEndpoint {
	if d.ce != nil {
		return d.ce
	}

	cl := &CollectionEndpoint{}
	cl.client = d.client.Clone()
	cl.client.BaseUrl().Path += CollectionPath
	cl.database = d
	d.ce = cl

	return d.ce
}

//DocumentEndPoint gets the document endpoint
func (d *Database) DocumentEndpoint() *DocumentEndpoint {

	if d.de != nil {
		return d.de
	}

	doc := &DocumentEndpoint{}
	doc.client = d.client.Clone()
	doc.client.BaseUrl().Path += DocumentPath
	doc.database = d
	d.de = doc

	return d.de
}

//EdgeEndPoint gets the document endpoint
func (d *Database) EdgeEndpoint() *EdgeEndpoint {
	if d.ee != nil {
		return d.ee
	}

	edge := &EdgeEndpoint{}
	edge.client = d.client.Clone()
	edge.client.BaseUrl().Path += EdgePath
	edge.database = d
	d.ee = edge

	return d.ee
}

//SimpleEndPoint gets the simple endpoint for simple queries
func (d *Database) SimpleEndpoint() *SimpleEndpoint {
	if d.se != nil {
		return d.se
	}
	s := &SimpleEndpoint{}
	s.client = d.client.Clone()
	s.client.BaseUrl().Path += SimplePath
	s.database = d
	d.se = s

	return d.se
}

//CursorEndpoint gets the cursor endpoint
func (d *Database) CursorEndpoint() *CursorEndpoint {
	if d.curE != nil {
		return d.curE
	}

	c := &CursorEndpoint{}
	c.client = d.client.Clone()
	c.client.BaseUrl().Path += CursorPath
	c.database = d
	d.curE = c

	return d.curE
}

//Connection returns connection associated with this database.
//It should be non-nil
func (d *Database) Connection() *Connection {
	return d.connection
}

//Get -> GET on /_api/database
func (d *Database) Get() ([]string, error) {
	var result struct {
		Result []string `json:"result"`
	}
	var errorResult = ArangoError{}

	h, err := d.client.Get(&gr.Params{
		Path: DatabasePath,
		UnmarshalMap: gr.UnmarshalMap{
			http.StatusOK:         &result,
			http.StatusBadRequest: &errorResult,
			http.StatusForbidden:  &errorResult,
		},
	})

	if err != nil {
		return nil, err
	}

	if h.StatusCode != http.StatusOK {
		return nil, errorResult
	}

	return result.Result, nil

}

//GetUser -> GET on /_api/database/user
func (d *Database) GetUser() ([]string, error) {

	var result struct {
		Result []string `json:"result"`
	}

	var errorResult = ArangoError{}

	h, err := d.client.Get(&gr.Params{
		Path: DatabasePath + "/user",
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

	return result.Result, nil
}

type DatabaseDescriptor struct {
	Name     string `json:"name"`
	Id       string `json:"id"`
	Path     string `json:"path"`
	IsSystem bool   `json:"isSystem"`
}

//GetCurrent -> GET on /_api/database/current
func (d *Database) GetCurrent() (*DatabaseDescriptor, error) {

	var result struct {
		Result *DatabaseDescriptor `json:"result"`
	}
	var errorResult = ArangoError{}

	h, err := d.client.Get(&gr.Params{
		Path: DatabasePath + "/current",
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

	return result.Result, nil
}

//PostDatabaseOptions are options when using the PostDatabase method
type PostDatabaseOptions struct {
	Name  string `json:"name"`
	Users []User `json:"users,omitempty"`
}

//User can be used when posting a new database. It outlines the users
//that can access the database.
type User struct {
	Username string      `json:"username"`
	Passwd   string      `json:"passwd"`
	Active   bool        `json:"active"`
	Extra    interface{} `json:"extra"`
}

//Post -> POST on /_api/database
func (d *Database) Post(name string, opts *PostDatabaseOptions) error {

	var errorResult = ArangoError{}
	if opts == nil {
		opts = new(PostDatabaseOptions)
	}
	opts.Name = name

	h, err := d.client.Post(&gr.Params{
		Path: DatabasePath,
		Body: opts,
		UnmarshalMap: gr.UnmarshalMap{
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
			http.StatusConflict:   &errorResult,
			http.StatusForbidden:  &errorResult,
		},
	})

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusCreated {
		return errorResult
	}

	return nil
}

//Delete -> DELETE on /_api/database/{name}
func (d *Database) Delete(name string) error {

	var errorResult = ArangoError{}

	h, err := d.client.Delete(&gr.Params{
		Path: DatabasePath + "/" + name,
		UnmarshalMap: gr.UnmarshalMap{
			http.StatusBadRequest: &errorResult,
			http.StatusNotFound:   &errorResult,
		},
	})

	if err != nil {
		return err
	}

	if h.StatusCode != http.StatusOK {
		return errorResult
	}

	return nil
}
