package arango

import (
	"fmt"
	"net/url"
	"testing"
)

func setupConnection() *Connection {
	u, _ := url.Parse("http://root@localhost:8529")
	c, err := NewConnection(u)
	if err != nil {
		fmt.Print(err)
	}
	return c
}

func getDatabase(name string) *Database {
	c := setupConnection()
	return c.Database(name)
}

func getCE(database string) *CollectionEndpoint {
	var db = getDatabase(database)
	return db.CollectionEndpoint()
}

func getDE(database string) *DocumentEndpoint {
	return getDatabase(database).DocumentEndpoint()
}

func createTestCollection() (*CollectionEndpoint, string) {
	var ce = getCE("_system")
	opts := DefaultPostCollectionOptions()
	opts.Name = "test"

	ce.PostCollection(opts)

	return ce, opts.Name
}

func createTestEdgeCollection() (*CollectionEndpoint, string) {
	var ce = getCE("_system")
	opts := DefaultPostCollectionOptions()
	opts.Type = EDGE_COLLECTION
	opts.Name = "test-edge"

	ce.PostCollection(opts)

	return ce, opts.Name
}

/*
func getEE(database string) *EdgeEndpoint {
	return getDatabase(database).EdgeEndpoint()
}
*/
/*
func getSE(database string) *SimpleEndpoint {
	return getDatabase(database).SimpleEndpoint()
}
*/

func verifyError(err error, t *testing.T, code int, message string) {
	if err == nil {
		t.Fatal(message)
	}

	ae, ok := err.(ArangoError)

	if !ok {
		t.Fatalf("Expected an ArangoError to be returned. %#v", err)
	}

	if !ae.IsError {
		t.Fatalf("Actual ae.IsError() == %t, Expected true. %#v, Message = %s", ae.IsError, ae, message)
	}

	if ae.Code != code {
		t.Fatalf("Actual ae.Code() == %d, Expected %d. %#v, Message = %s", ae.Code, code, ae, message)
	}
}
