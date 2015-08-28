package arango

import (
	"net/url"
	"testing"
)

func TestNilUrlFails(t *testing.T) {
	var err error
	_, err = NewConnection(nil)

	if err == nil {
		t.Fatal("Expected error when creating a new connection with a nil url.")
	}
}

func TestBadHostError(t *testing.T) {
	var err error
	u, err := url.Parse("http://badhost:8529")
	c, err := NewConnection(u)

	_, err = c.Version(false)
	if err == nil {
		t.Fatal("Expected error with bad host.")
	}
}

func TestGetVersion(t *testing.T) {
	var err error
	c := setupConnection()

	v, err := c.Version(false)

	if err != nil {
		t.Fatal("Could not get version: ", err)
	}

	if v.Server != "arango" {
		t.Fatal("Unexpected server value: ", v.Server)
	}

	if v.Version != "2.6.5" {
		t.Fatal("Unexpected version value: ", v.Version)
	}

	if v.Details != nil || len(v.Details) > 0 {
		t.Fatal("Unexpected details when none were requested.", v.Details)
	}

}

func TestGetVersionWithDetails(t *testing.T) {
	var err error
	c := setupConnection()

	v, err := c.Version(true)

	if err != nil {
		t.Fatal("Could not get version: ", err)
	}

	if v.Details == nil || len(v.Details) < 1 {
		t.Fatal("Unable to fetch details.", v.Details)
	}
}
