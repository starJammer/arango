package arango

import (
    
    "testing"
)

var (
    db *Database
)

func setup() {
    var err error
    db, err = Conn( "http://root@localhost:8529" )
    if err != nil {
        panic( err )
    }

    db.DropDatabase( "testing" )
    err = db.CreateDatabase( "testing", nil, users )

    if err != nil {
        panic( err )
    }
    db, err = db.UseDatabase( "testing" )
    if err != nil {
        panic( err )
    }
}

func teardown(){
    var err error
    db, err = db.UseDatabase( "_system" )
    if err != nil {
        panic( err )
    }
    err = db.DropDatabase( "testing" )
    if err != nil {
        panic( err )
    }
}

type DummyDocument struct{
    Hi string
}

type DummyFullDocument struct {
    DocumentImplementation
    Hi string
}

func TestSavingAndRetrievingDocument( t *testing.T ){

    setup()
    defer teardown()

    c, err := db.CreateDocumentCollection( "testing" )

    err = c.Save( &DummyDocument{ Hi : "Hello World" } )

    if err != nil {
        t.Fatal( err )
    }

    fulld := &DummyFullDocument{
        Hi : "Hello World",
    }

    err = c.Save( fulld )

    if err != nil {
        t.Fatal( err )
    }

    if fulld.Id() == "" {
        t.Fatal( "Expected id to be populated in document after a save." )
    }

    //Find it via full id
    ret := &DummyFullDocument{}
    err = c.Document( fulld.Id(), ret )

    if err != nil {
        t.Fatal( err )
    }

    if ret.Hi != "Hello World" {
        t.Fatal( "Expected to have the value for the document correctly fetched.")
    }

    if ret.Id() != fulld.Id() {
        t.Fatal( "Expected to have the ids for documents be equal since they are the same document.")
    }

    if ret.Key() != fulld.Key() {
        t.Fatal( "Expected to have the keys for documents be equal since they are the same document.")
    }

    if ret.Rev() != fulld.Rev() {
        t.Fatal( "Expected to have the revs for documents be equal since they are the same document.")
    }

    //Find it using a Document struct
    ret = &DummyFullDocument{}
    err = c.Document( fulld, ret )

    if err != nil {
        t.Fatal( err )
    }

    if ret.Hi != "Hello World" {
        t.Fatal( "Expected to have the value for the document correctly fetched.")
    }

    if ret.Id() != fulld.Id() {
        t.Fatal( "Expected to have the ids for documents be equal since they are the same document.")
    }

    if ret.Key() != fulld.Key() {
        t.Fatal( "Expected to have the keys for documents be equal since they are the same document.")
    }

    if ret.Rev() != fulld.Rev() {
        t.Fatal( "Expected to have the revs for documents be equal since they are the same document.")
    }

    //Find it by key
    ret = &DummyFullDocument{}
    err = c.Document( fulld.Key(), ret )

    if err != nil {
        t.Fatal( err )
    }

    if ret.Hi != "Hello World" {
        t.Fatal( "Expected to have the value for the document correctly fetched.")
    }

    if ret.Id() != fulld.Id() {
        t.Fatal( "Expected to have the ids for documents be equal since they are the same document.")
    }

    if ret.Key() != fulld.Key() {
        t.Fatal( "Expected to have the keys for documents be equal since they are the same document.")
    }

    if ret.Rev() != fulld.Rev() {
        t.Fatal( "Expected to have the revs for documents be equal since they are the same document.")
    }

    //Try a cross collection update
    err = c.Document( "fake/"+fulld.Key(), ret )

    if err == nil {
        t.Fatal( "Expected a cross collection error.")
    }

    fulld.SetId( "fake/" + fulld.Key() )

    err = c.Document( fulld, ret )

    if err == nil {
        t.Fatal( "Expected a cross collection error.")
    }

}
