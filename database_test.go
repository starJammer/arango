package arango

import (
	"net/url"
	"testing"
)

func TestGetDatabase(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)

	var db Database = c.Database("_system")

	if db.Name() != "_system" {
		t.Fatal("Database name incorrect: ", db.Name())
	}

}

func TestGetDatabases(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)

	var db Database = c.Database("_system")

	dbs, err := db.Get()

	if err != nil {
		t.Fatal("Could not fetch all the databases: ", err)
	}

	if len(dbs) < 1 {
		t.Fatal("Expected to at least have the _system database but don't.")
	}

	userDbs, err := db.GetUser()

	if err != nil {
		t.Fatal("Could not fetch user databases: ", err)
	}

	if len(userDbs) != len(dbs) {
		t.Fatalf("User dbs not the same as regular dbs - user: %s, regular: %s", userDbs, dbs)
	}
}

func TestGetCurrent(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)

	var db Database = c.Database("_system")

	current, err := db.GetCurrent()

	if err != nil {
		t.Fatal("Could not fetch current db info:", err)
	}

	if current.Name() != "_system" {
		t.Fatal("Unexpected name: ", current.Name())
	}

	if current.Id() == "" {
		t.Fatal("Unexpected id: ", current.Id())
	}

	if current.Path() == "" {
		t.Fatal("Unexpected path: ", current.Path())
	}

	if !current.IsSystem() {
		t.Fatalf("Unexpected isSystem: %t", current.IsSystem())
	}
}

func TestCreateDeleteDatabase(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)

	var db Database = c.Database("_system")
	var err error

	err = db.Post("", nil)

	if err == nil {
		t.Fatal("Expected error creating a database with no name.")
	}

	err = db.Post("_system", nil)

	if err == nil {
		t.Fatal("Expected error when creating a database that exists already: ")
	}

	err = db.Post("test", nil)

	if err != nil {
		t.Fatal("Unexpected error when creating new database: ", err)
	}

	err = db.Delete("")
	if err == nil {
		t.Fatal("Expected error when deleting database with no name.")
	}

	err = db.Delete("test")

	if err != nil {
		t.Fatal("Unexpected error when deleting a database: ", err)
	}
}
