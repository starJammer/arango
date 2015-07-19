package arango

import (
	gr "github.com/starJammer/grestclient"
)

type collection struct {
	name     string
	client   gr.Client
	database *database
}

func (c *collection) Name() string {
	return c.name
}

func (c *collection) Database() Database {
	return c.database
}

func (c *collection) Get() error {
	return nil
}

func (c *collection) GetProperties() error {
	return nil
}

func (c *collection) Delete() error {

	var errorResult = &arangoError{}

	h, err := c.client.Delete("", nil, nil, errorResult)

	if err != nil {
		return err
	}

	if h.StatusCode != 200 {
		return errorResult
	}

	return nil

}
