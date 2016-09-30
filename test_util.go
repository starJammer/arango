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
		t.Fatalf("Actual ae.IsError() == %t, Expected true. ArangoError = %s, Message = %s", ae.IsError, ae, message)
	}

	if ae.Code != code {
		t.Fatalf("Actual ae.Code() == %d, Expected %d. ArangoError = %s, Message = %s", ae.Code, code, ae, message)
	}
}
