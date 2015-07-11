package old

import (
	"testing"
	//"log"
)

var (
	users = []User{
		User{Username: "alice", Passwd: "hi", Active: true, Extra: struct{}{}},
		User{Username: "bob", Passwd: "hi", Active: false, Extra: struct{}{}},
		User{Username: "charlie", Passwd: "hi", Active: true, Extra: struct{ Extra string }{Extra: "hi"}},
		User{Username: "root", Passwd: "", Active: true, Extra: struct{ Extra string }{Extra: "hi"}},
	}
)

func TestDatabaseCreateUseDropMethods(t *testing.T) {
	var e ArangoError

	db, err := Conn("http://root@localhost:8529")
	if err != nil {
		t.Fatal(err)
	}

	err = db.DropDatabase("testing")

	if err != nil {
		e, _ := err.(ArangoError)
		if e.Code != 404 && e.Code == 200 {
			t.Log("Database exists but was expecting it not to.", err)
		}
	}

	err = db.CreateDatabase("testing", nil, users)

	if err != nil {
		t.Fatal(err)
	}

	//Trying to create testing again should fail
	err = db.CreateDatabase("testing", nil, nil)

	if err == nil {
		t.Fatalf("Creating a database twice should have produced an error.")
	}
	_, ok := err.(ArangoError)

	if !ok {
		t.Fatalf("Did not get the expected error type : ArangoError")
	}

	db, err = db.UseDatabase("testing")

	if err != nil {
		t.Fatal(err)
	}

	if db.Name() != "testing" {
		t.Fatal("Failed to switch to the database :", "testing")
	}

	if db.IsSystem() {
		t.Fatal("The database we switched to SHOULD NOT be a system database but is : ", "testing")
	}

	err = db.DropDatabase("testing")

	if err == nil {
		t.Fatal("Dropping the database we are in should fail. Drop can only be called from _system")
	}

	e, ok = err.(ArangoError)

	if !ok || !e.IsError {
		t.Error(err)
		t.Fatal("Dropping the database we are in should fail. Drop can only be called from _system")
	}

	db, err = db.UseDatabase("_system")
	err = db.DropDatabase("testing")

	if err != nil {
		t.Fatal("Dropping the database should work from the _system database.", err)
	}

	db, err = db.UseDatabase("testing")

	if err == nil {
		t.Fatal("Expected to get an error for using a database that doesn't exist.")
	}

}

func TestDatabaseCreateUseDropMethodsUnixConnection(t *testing.T) {
	var e ArangoError

	db, err := Conn("unix://root@/tmp/arangod.soc")
	if err != nil {
		t.Fatal(err)
	}

	err = db.DropDatabase("testing")

	if err != nil {
		e, _ := err.(ArangoError)
		if e.Code != 404 && e.Code == 200 {
			t.Log("Database exists but was expecting it not to.", err)
		}
	}

	err = db.CreateDatabase("testing", nil, users)

	if err != nil {
		t.Fatal(err)
	}

	//Trying to create testing again should fail
	err = db.CreateDatabase("testing", nil, nil)

	if err == nil {
		t.Fatalf("Creating a database twice should have produced an error.")
	}
	_, ok := err.(ArangoError)

	if !ok {
		t.Fatalf("Did not get the expected error type : ArangoError")
	}

	db, err = db.UseDatabase("testing")

	if err != nil {
		t.Fatal(err)
	}

	if db.Name() != "testing" {
		t.Fatal("Failed to switch to the database :", "testing")
	}

	if db.IsSystem() {
		t.Fatal("The database we switched to SHOULD NOT be a system database but is : ", "testing")
	}

	err = db.DropDatabase("testing")

	if err == nil {
		t.Fatal("Dropping the database we are in should fail. Drop can only be called from _system")
	}

	e, ok = err.(ArangoError)

	if !ok || !e.IsError {
		t.Error(err)
		t.Fatal("Dropping the database we are in should fail. Drop can only be called from _system")
	}

	db, err = db.UseDatabase("_system")
	err = db.DropDatabase("testing")

	if err != nil {
		t.Fatal("Dropping the database should work from the _system database.", err)
	}

	db, err = db.UseDatabase("testing")

	if err == nil {
		t.Fatal("Expected to get an error for using a database that doesn't exist.")
	}

}

func TestDatabaseCollectionMethods(t *testing.T) {

	db, err := Conn("http://root@localhost:8529")
	if err != nil {
		t.Fatal(err)
	}

	err = db.DropDatabase("testing")

	if err != nil {
		e, _ := err.(ArangoError)
		if e.Code != 404 && e.Code == 200 {
			t.Log("Database exists but was expecting it not to.", err)
		}
	}

	db.CreateDatabase("testing", nil, users)

	db, err = db.UseDatabase("testing")

	c, err := db.CreateDocumentCollection("testing")

	if err != nil {
		t.Fatalf("Was not expecting an error when creating testing collection:%s", err)
	}

	//Try fetching a collection that shouldn't exist
	_, err = db.Collection("testing_no_exist")

	if err == nil {
		t.Fatal("Expecting an error because the collection didn't exist.")
	}

	if c == nil {
		t.Fatal("Expecting the collection returned to not be nil.")
	}

	if c.Id() == "" {
		t.Fatal("Collection should have an Id associated with it.")
	}

	if c.Status() == 0 {
		t.Fatal("Collection should have a status other than 0 after creation.")
	}

	if c.Name() != "testing" {
		t.Fatal("Collection doesn't have expected name.")
	}

	if c.Type() != DOCUMENT_COLLECTION {
		t.Fatalf("Collection should be document (%d) type but got something else (%d).", DOCUMENT_COLLECTION, c.Type())
	}

	if c.IsSystem() {
		t.Fatal("Collection should not be a system collection.")
	}

	//Fetch the database using the Collection method
	c, err = db.Collection("testing")

	if err != nil {
		t.Fatal("Got an unexpected error when getting the testing collection.")
	}

	if c == nil {
		t.Fatal("Expecting the collection returned to not be nil.")
	}

	if c.Id() == "" {
		t.Fatal("Collection should have an Id associated with it.")
	}

	if c.Status() == 0 {
		t.Fatal("Collection should have a status other than 0 after creation.")
	}

	if c.Name() != "testing" {
		t.Fatal("Collection doesn't have expected name.")
	}

	if c.Type() != DOCUMENT_COLLECTION {
		t.Fatalf("Collection should be document (%d) type but got something else (%d).", DOCUMENT_COLLECTION, c.Type())
	}

	if c.IsSystem() {
		t.Fatal("Collection should not be a system collection.")
	}

	err = c.Properties() //fetch properties

	if err != nil {
		t.Fatal("Got an error when fetching properties.", err)
	}

	if c.KeyOptions() == nil {
		t.Fatal("KeyOptions for the collection should not be nil.")
	}

	if c.KeyOptions().Type == "" {
		t.Fatal("KeyOptions.Type for the collection should not be blank.")
	}

	if c.JournalSize() == 0 {
		t.Fatal("JournalSize for the collection should not be 0.")
	}

	err = db.DropCollection("testing_no_exist")

	if err == nil {
		t.Fatalf("Expected an error when dropping a non-existent collection.")
	}

	err = db.DropCollection(c.Name())

	if err != nil {
		t.Fatalf("Could not drop the collection: %+v", err)
	}

	//Clean up everything
	db, err = db.UseDatabase("_system")
	err = db.DropDatabase("testing")

	if err != nil {
		t.Fatal("Dropping the database should work from the _system database.")
	}

}

func TestDatabaseDocumentMethods(t *testing.T) {

	setup()
	defer teardown()

	db := db

	type basic struct {
		DocumentImplementation
		Hi string
	}

	a := &basic{Hi: "hey"}
	err := db.SaveDocumentWithOptions(a, &SaveOptions{})

	if err == nil {
		t.Fatal("Expected an error because no collection was specified.")
	}

	err = db.SaveDocumentWithOptions(a, &SaveOptions{Collection: "testing"})

	if err == nil {
		t.Fatal("Expected an error because the collection didn't exist.")
	}

	err = db.SaveDocumentWithOptions(a, &SaveOptions{Collection: "testing", CreateCollection: true})

	if err != nil {
		t.Fatal("Did not expect error when saving to collection:", err)
	}

	oRev := a.Rev()
	oKey := a.Key()
	oId := a.Id()
	a.SetRev("")
	a.SetKey("")
	a.SetId("")
	//Try it again and we still should not receive an error because the collection already exists
	//Just a double check for arango working as expected, not so much our code
	err = db.SaveDocumentWithOptions(a, &SaveOptions{Collection: "testing"})

	if err != nil {
		t.Fatal("Did not expect error when saving to collection:", err)
	}

	if oRev == a.Rev() || oKey == a.Key() || oId == a.Id() {
		t.Fatal("Did not expect Arango to return the same revision for a new object.")
	}

	err = db.SaveDocumentWithOptions(&basic{Hi: "hey"}, &SaveOptions{Collection: "testing", WaitForSync: true})

	if err != nil {
		t.Fatal("Did not expect error when saving to collection:", err)
	}

	//Now test fetching
	var b basic

	err = db.Document(a.Id(), &b)

	if err != nil {
		t.Fatal("Did not expect an error when fetching by id.", err)
	}

	if b.Hi != "hey" {
		t.Fatal("Expected property to be \"hey\"")
	}

	b.Hi = "Testing second fetch"
	err = db.Document(a.Id(), &b)

	if b.Hi != "hey" {
		t.Fatal("Expected property to be set back to \"hey\"")
	}

	//Test updating a document
	var newString = "new string not there before."
	a.Hi = newString
	err = db.UpdateDocumentWithOptions(a.Id(), a, nil)

	if err != nil {
		t.Fatal("Did not expect update on document to fail.", err)
	}

	if a.Rev() == b.Rev() {
		t.Fatal("The revision of the old document and the new document should not be the same.")
	}

	//Try to update with the older Revision as a condition. Should fail.
	upOpts := DefaultUpdateOptions()
	upOpts.Rev = b.Rev()
	err = db.UpdateDocumentWithOptions(b.Id(), &b, upOpts)

	if err == nil {
		t.Fatal("Expected an error because the revision of the document does not match the one in the database anymore.")
	}

	var c = &struct {
		DocumentImplementation
		C string `json:"c"`
	}{C: "new property"}

	var d = &struct {
		DocumentImplementation
		Hi string
		C  string `json:"c"`
	}{}

	//Test adding a property to a document
	err = db.UpdateDocumentWithOptions(a.Id(), c, nil)

	if err != nil {
		t.Fatal("Did not expect update on document to fail.", err)
	}

	err = db.DocumentWithOptions(a.Id(), d, nil)

	if err != nil {
		t.Fatal("Did not expect fetching document to fail.", err)
	}

	if d.Hi != a.Hi || d.C != c.C {
		t.Fatalf("Structs d did not have the correct properties after an update. \n(a, %+v)\n(c, %+v)\n(d, %+v)", a, c, d)
	}

	//Test deleting a document
	err = db.DeleteDocumentWithOptions(a.Id(), nil)

	if err != nil {
		t.Fatal("Did not expect delete on document to fail.", err)
	}

	err = db.DocumentWithOptions(a.Id(), d, nil)

	if err == nil {
		t.Fatal("Expected an error because document should not exist anymore.", err)
	}

}

func TestDatabaseEdgeMethods(t *testing.T) {

	setup()
	defer teardown()

	db := db

	type basic struct {
		DocumentImplementation
	}

	type basicEdge struct {
		EdgeImplementation
		Field string `json:"name"`
	}

	type basicEdge2 struct {
		EdgeImplementation
		Field2 string `json:"name2"`
	}

	t1, err := db.CreateDocumentCollection("thing1")
	t2, err := db.CreateDocumentCollection("thing2")
	edges, err := db.CreateEdgeCollection("edges")

	if err != nil {
		t.Fatal("Could not create the edge collection!")
	}

	var a = new(basic)
	var b = new(basic)
	var e = new(basicEdge)
	var f = new(basicEdge)

	e.Field = "testing for an edge"

	t1.Save(a)
	t2.Save(b)

	err = edges.SaveEdge(a, b, e)

	if err != nil {
		t.Fatal("An unexpected error occurred when saving the edge :", err)
	}

	if e.From() != a.Id() {
		t.Fatalf("Expected the FROM in the edge to be correct. (%s, %s)", e.From(), a.Id())
	}

	if e.To() != b.Id() {
		t.Fatalf("Expected the TO in the edge to be correct. (%s, %s)", e.To(), b.Id())
	}

	err = db.Edge(e, f)

	if err != nil {
		t.Fatal("Could not fetch the edge we just created.")
	}

	if f.Id() != e.Id() {
		t.Fatal("The _id values of the edges are not equal.")
	}

	if f.Rev() != e.Rev() {
		t.Fatal("The _rev values of the edges are not equal.")
	}

	if f.From() != e.From() {
		t.Fatal("The FROM values of the edges are not equal.")
	}

	if f.To() != e.To() {
		t.Fatal("The TO values of the edges are not equal.")
	}

	e.Field = "something new"

	err = db.UpdateEdgeWithOptions(e.Id(), e, nil)

	if err != nil {
		t.Fatal("Could not update the edge.")
	}

	if e.Rev() == f.Rev() {
		t.Fatal("Revision numbers should not be equal anymore.")
	}

	db.Edge(e.Id(), f)

	if e.Field != f.Field {
		t.Fatal("Fields should now be equal but they're not.")
	}

	var replacement = new(basicEdge2)
	replacement.Field2 = "yeah"

	err = db.ReplaceEdgeWithOptions(e.Id(), replacement, nil)

	if err != nil {
		t.Fatal("Replacing the edge did not work.")
	}

	e = new(basicEdge)
	db.Edge(replacement.Id(), e)

	if e.Field != "" {
		t.Fatal("The old Field was not deleted during the replacement.")
	}

}
