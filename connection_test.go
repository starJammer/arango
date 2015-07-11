package arango

import (
	"net/url"
	"testing"
)

func TestArangoErrorMeetsInterface(t *testing.T) {
	var _ ArangoError = &arangoError{}
}

func TestConnectionMeetsInterface(t *testing.T) {
	var err error
	u, err := url.Parse("http://localhost:8529")
	c, err := NewConnection(u)

	if err != nil {
		t.Fatal("Could not create a connection: ", err)
	}

	if _, ok := c.(HasGrestClient); !ok {
		t.Fatal("Expected this implementation to use grestclient.")
	}
}

func TestGetVersion(t *testing.T) {
	var err error
	u, err := url.Parse("http://root@localhost:8529")
	c, err := NewConnection(u)

	v, err := c.Version(false)

	if err != nil {
		t.Fatal("Could not get version: ", err)
	}

	if v.Server() != "arango" {
		t.Fatal("Unexpected server value: ", v.Server())
	}

	if v.Version() != "2.6.2" {
		t.Fatal("Unexpected version value: ", v.Version())
	}

	if v.Details() != nil || len(v.Details()) > 0 {
		t.Fatal("Unexpected details when none were requested.", v.Details())
	}

	v, err = c.Version(true)

	if v.Details() == nil || len(v.Details()) < 1 {
		t.Fatal("Unable to fetch details.", v.Details())
	}
}

func TestGetDatabase(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)

	var db Database = c.Database("_system")

	if db.Name() != "_system" {
		t.Fatal("Database name incorrect: ", db.Name())
	}

}
