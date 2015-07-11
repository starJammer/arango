package old

import (
	"testing"
)

func TestConnectionSuccessful(t *testing.T) {
	db, err := Conn("http://root@localhost:8529")

	if err != nil {
		t.Fatal(err)
	}

	if db.Name() != "_system" {
		t.Error("Did not get default database _system.")
	}

	if !db.IsSystem() {
		t.Error("Expected _system to have IsSystem = true")
	}

	if db.Path() == "" {
		t.Error("Path of the database was not set. A path is expected.")
	}

	if db.Id() == "" {
		t.Error("Id of the database was not set. An id is expected.")
	}

}

func TestSslConnectionSuccessful(t *testing.T) {
	AllowBadSslCerts = true
	_, err := Conn("https://root@localhost:8530")
	AllowBadSslCerts = false
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnixSocketConnectionSuccessful(t *testing.T) {
	_, err := Conn("unix://root@/tmp/arangod.soc")
	if err != nil {
		t.Fatal(err)
	}
}

func TestBadPathUnixSocketConnectionSuccessful(t *testing.T) {
	_, err := Conn("unix://root@/tmp/fakearangod.soc")
	if err == nil {
		t.Fatal("Expected error when connectiong to unix:///tmp/fakearangod.soc")
	}
}

func TestConnectionFailure(t *testing.T) {
	_, err := Conn("http://root@localhost:9999")

	if err == nil {
		t.Fatal("Expected error when connectiong to http://localhost:9999")
	}

	if _, ok := err.(ArangoError); !ok {
		t.Fatalf("Expected ArangoError but got something else (%T, %v)", err, err)
	}
}

func TestBadUser(t *testing.T) {
	db, err := Conn("http://roo@localhost:8529")

	if err == nil {
		t.Fatal("Expected error when connectiong to http://localhost:8529")
	}

	if _, ok := err.(ArangoError); !ok {
		t.Fatalf("Expected ArangoError but got something else (%T, %v)", err, err)
	}

	if db != nil {
		t.Fatalf("Expected database to be nil but got something back : (%T, %v)", db, db)
	}
}

func TestUsingDatabaseName(t *testing.T) {
	db, err := ConnDb("http://root@localhost:8529", "_system")

	if err != nil {
		t.Fatal(err)
	}

	if db.Name() != "_system" {
		t.Error("Did not get default database _system.")
	}

}

func TestUsingDatabaseNameUnixConn(t *testing.T) {
	db, err := ConnDb("unix://root@/tmp/arangod.soc", "_system")

	if err != nil {
		t.Fatal(err)
	}

	if db.Name() != "_system" {
		t.Error("Did not get default database _system.")
	}

}

func TestUsingDatabaseNameAndUserCreds(t *testing.T) {
	db, err := ConnDbUserPassword("http://localhost:8529", "_system", "root", "")

	if err != nil {
		t.Fatal(err)
	}

	if db == nil {
		t.Fatal("Expected to get a database back but got nil.")
	}

	if db.Name() != "_system" {
		t.Error("Did not get default database _system.")
	}

}

func TestFailUsingDatabaseNameAndUserCreds(t *testing.T) {
	db, err := ConnDbUserPassword("http://localhost:8529", "_system", "roo", "")

	if err == nil {
		t.Fatal("Expected error when connectiong to http://localhost:8529")
	}

	if _, ok := err.(ArangoError); !ok {
		t.Fatalf("Expected ArangoError but got something else (%T, %v)", err, err)
	}

	if db != nil {
		t.Fatalf("Expected database to be nil but got something back : (%T, %v)", db, db)
	}

}
