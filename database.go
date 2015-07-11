package arango

import (
	"fmt"
	gr "github.com/starJammer/grestclient"
)

type database struct {
	connection Connection
	client     gr.Client
}

func (d *database) Collection(name string) Collection {
	return nil
}

func (d *database) Connection() Connection {
	return d.connection
}

func (d *database) Get() ([]string, error) {

	return []string{}, nil
}

func (d *database) Delete(name string) error {
	return nil
}

func (c *connection) Database(name string) Database {
	db := &database{}
	db.connection = c
	db.client = c.client.Clone()
	db.client.BaseUrl().Path += fmt.Sprintf(Databasepath, name)
	return db
}
