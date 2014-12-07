package arango

import (
    "testing"
    //"log"
)

func TestDatabaseCreateUseDropMethods( t *testing.T ) {
    var e ArangoError

    db, err := Conn( "http://root@localhost:8529")
    if err != nil {
        t.Fatal( err )
    }

    err = db.DropDatabase( "testing" )

    if err != nil {
        e, _ := err.(ArangoError)
        if e.Code != 404 && e.Code == 200 {
            t.Log( "Database exists but was expecting it not to.", err )
        }
    }

    err = db.CreateDatabase( "testing", nil, nil )

    if err != nil { 
        t.Fatal( err )
    }

    //Trying to create testing again should fail
    err = db.CreateDatabase( "testing", nil, nil )

    if err == nil {
        t.Fatalf( "Creating a database twice should have produced an error." )
    }
    _, ok := err.(ArangoError)

    if !ok {
        t.Fatalf( "Did not get the expected error type : ArangoError" )
    }

    db, err = db.UseDatabase( "testing" )

    if err != nil {
        t.Fatal( err )
    }

    if db.Name() != "testing" {
        t.Fatal( "Failed to switch to the database :", "testing")
    }

    if db.IsSystem() {
        t.Fatal( "The database we switched to SHOULD NOT be a system database but is : ", "testing")
    }

    err = db.DropDatabase( "testing" )

    if err == nil {
        t.Fatal( "Dropping the database we are in should fail. Drop can only be called from _system" )
    }

    e, ok = err.(ArangoError)

    if !ok || !e.IsError {
        t.Error( err )
        t.Fatal( "Dropping the database we are in should fail. Drop can only be called from _system" )
    }

}
