package arango

import (
	"fmt"
	gr "github.com/starJammer/grestclient"
)

func (c *connection) Database(name string) Database {
	db := &database{}
	db.connection = c
	db.name = name
	db.client = c.client.Clone()
	db.client.BaseUrl().Path += fmt.Sprintf(Databasepath, name)

	return db
}

type database struct {
	connection Connection
	client     gr.Client
	name       string
}

func (d *database) Name() string {
	return d.name
}

func (d *database) Collection(name string) Collection {
	return nil
}

func (d *database) Connection() Connection {
	return d.connection
}

func (d *database) Get() ([]string, error) {
	var result struct {
		Result []string `json:"result"`
	}
	var errorResult = &arangoError{}

	h, err := d.client.Get(DatabaseEndPoint, nil, &result, errorResult)

	if err != nil {
		return nil, err
	}

	if h.StatusCode >= 400 {
		return nil, errorResult
	}

	return result.Result, nil

}

func (d *database) GetUser() ([]string, error) {

	var result struct {
		Result []string `json:"result"`
	}

	var errorResult = &arangoError{}

	h, err := d.client.Get(DatabaseEndPoint+"/user", nil, &result, errorResult)

	if err != nil {
		return nil, err
	}

	if h.StatusCode >= 400 {
		return nil, errorResult
	}

	return result.Result, nil
}

type currentResult struct {
	Namef     string `json:"name"`
	Idf       string `json:"id"`
	Pathf     string `json:"path"`
	IsSystemf bool   `json:"isSystem"`
}

func (cr *currentResult) Name() string {
	return cr.Namef
}

func (cr *currentResult) Id() string {
	return cr.Idf
}

func (cr *currentResult) Path() string {
	return cr.Pathf
}

func (cr *currentResult) IsSystem() bool {
	return cr.IsSystemf
}

func (d *database) GetCurrent() (CurrentResult, error) {

	var result struct {
		Result *currentResult `json:"result"`
	}
	var errorResult = &arangoError{}

	h, err := d.client.Get(DatabaseEndPoint+"/current", nil, &result, errorResult)

	if err != nil {
		return nil, err
	}

	if h.StatusCode >= 400 {
		return nil, errorResult
	}

	return result.Result, nil
}

func (d *database) Post(opts *PostDatabaseOptions) error {

	var errorResult = &arangoError{}

	h, err := d.client.Post(DatabaseEndPoint, nil, opts, nil, errorResult)

	if err != nil {
		return err
	}

	if h.StatusCode != 201 {
		return errorResult
	}

	return nil
}

func (d *database) Delete(name string) error {

	var errorResult = &arangoError{}

	h, err := d.client.Delete(
		DatabaseEndPoint+"/"+name,
		nil,
		nil, errorResult)

	if err != nil {
		return err
	}

	if h.StatusCode != 200 {
		return errorResult
	}

	return nil
}

func (d *database) GetCollections() {

}

func (d *database) PostCollection() {

}
