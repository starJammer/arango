package arango

import (
    "testing"
)

func TestConnectionSuccessful( t *testing.T ){
    db, err := Conn( "http://root@localhost:8529" )

    if err != nil {
        t.Fatal( err )
    }

    if db.Name() != "_system" {
        t.Error( "Did not get default database _system.")
    }

    if !db.IsSystem() {
        t.Error( "Expected _system to have IsSystem = true")
    }

    if db.Path() == "" {
        t.Error( "Path of the database was not set. A path is expected." )
    }

    if db.Id() == "" {
        t.Error( "Id of the database was not set. An id is expected." )
    }

}

func TestSslConnectionSuccessful( t *testing.T ){
    AllowBadSslCerts = true
    _, err := Conn( "https://root@localhost:8530" )
    AllowBadSslCerts = false
    if err != nil {
        t.Fatal( err )
    }
}

func TestConnectionFailure( t *testing.T ){
    _, err := Conn( "http://root@localhost:9999" )

    if err == nil {
        t.Fatal( "Expected error when connectiong to http://localhost:9999" )
    }

    if _, ok := err.(ArangoError); !ok {
        t.Fatalf( "Expected ArangoError but got something else (%T, %v)", err )
    }
}

func TestBadUser( t *testing.T ){
    _, err := Conn( "http://roo@localhost:8529" )

    if err == nil {
        t.Fatal( "Expected error when connectiong to http://localhost:9999" )
    }

    if _, ok := err.(ArangoError); !ok {
        t.Fatalf( "Expected ArangoError but got something else (%T, %v)", err )
    }
}
