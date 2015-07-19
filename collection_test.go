package arango

import (
	"net/url"
	"testing"
)

func TestGetCollection(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")

	opts := DefaultCollectionOptions()
	opts.Name = "test"

}
