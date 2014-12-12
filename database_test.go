package arango

import (
    "testing"
    //"log"
)

var (
 users = []User{ 
        User{ Username : "alice", Passwd : "hi", Active : true, Extra : struct{}{}  },
        User{ Username : "bob", Passwd : "hi", Active : false, Extra : struct{}{}  },
        User{ Username : "charlie", Passwd : "hi", Active : true, Extra : struct{ Extra string }{ Extra : "hi"}  },
        User{ Username : "root", Passwd : "", Active : true, Extra : struct{ Extra string }{ Extra : "hi"}  },
    }
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

    err = db.CreateDatabase( "testing", nil, users )

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

    db, err = db.UseDatabase( "_system" )
    err = db.DropDatabase( "testing" )

    if err != nil {
        t.Fatal( "Dropping the database should work from the _system database.", err )
    }

    db, err = db.UseDatabase( "testing" )

    if err == nil {
        t.Fatal( "Expected to get an error for using a database that doesn't exist." )
    }

}

func TestDatabaseCollectionMethods( t *testing.T ) {

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

    db.CreateDatabase( "testing", nil, users )

    db, err = db.UseDatabase( "testing" )

    c, err := db.CreateDocumentCollection( "testing" )

    if err != nil {
        t.Fatalf( "Was not expecting an error when creating testing collection:%s", err)
    }

    //Try fetching a collection that shouldn't exist
    _, err = db.Collection( "testing_no_exist" )

    if err == nil {
        t.Fatal( "Expecting an error because the collection didn't exist." )
    }

    if c == nil {
        t.Fatal( "Expecting the collection returned to not be nil.")
    }

    if c.Id() == "" {
        t.Fatal( "Collection should have an Id associated with it." )
    }

    if c.Status() == 0 {
        t.Fatal( "Collection should have a status other than 0 after creation." )
    }

    if c.Name() != "testing" {
        t.Fatal( "Collection doesn't have expected name.")
    }

    if c.Type() != DOCUMENT_COLLECTION {
        t.Fatalf( "Collection should be document (%d) type but got something else (%d).", DOCUMENT_COLLECTION, c.Type() )
    }

    if c.IsSystem() {
        t.Fatal( "Collection should not be a system collection." )
    }

    //Fetch the database using the Collection method
    c, err = db.Collection( "testing" )

    if err != nil {
        t.Fatal( "Got an unexpected error when getting the testing collection.")
    }

    if c == nil {
        t.Fatal( "Expecting the collection returned to not be nil.")
    }

    if c.Id() == "" {
        t.Fatal( "Collection should have an Id associated with it." )
    }

    if c.Status() == 0 {
        t.Fatal( "Collection should have a status other than 0 after creation." )
    }

    if c.Name() != "testing" {
        t.Fatal( "Collection doesn't have expected name.")
    }

    if c.Type() != DOCUMENT_COLLECTION {
        t.Fatalf( "Collection should be document (%d) type but got something else (%d).", DOCUMENT_COLLECTION, c.Type() )
    }

    if c.IsSystem() {
        t.Fatal( "Collection should not be a system collection." )
    }

    err = c.Properties() //fetch properties

    if err != nil {
        t.Fatal( "Got an error when fetching properties.", err )
    }

    if c.KeyOptions() == nil {
        t.Fatal( "KeyOptions for the collection should not be nil.")
    }

    if c.KeyOptions().Type == "" {
        t.Fatal( "KeyOptions.Type for the collection should not be blank.")
    }

    if c.JournalSize() == 0 {
        t.Fatal( "JournalSize for the collection should not be 0.")
    }

    err = db.DropCollection( "testing_no_exist" )

    if err == nil {
        t.Fatalf( "Expected an error when dropping a non-existent collection.")
    }

    err = db.DropCollection( c.Name() )

    if err != nil {
        t.Fatalf( "Could not drop the collection: %+v", err )
    }

    //Clean up everything
    db, err = db.UseDatabase( "_system" )
    err = db.DropDatabase( "testing" )

    if err != nil {
        t.Fatal( "Dropping the database should work from the _system database." )
    }

}
