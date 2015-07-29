package arango

import (
	"net/url"
	"testing"
)

func TestPostGetDocument(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")
	var collEnd = db.CollectionEndpoint()
	var docEnd = db.DocumentEndpoint()

	docs := DefaultPostCollectionOptions()
	docs.Name = "test"

	err := collEnd.PostCollection(docs.Name, nil)
	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer collEnd.Delete(docs.Name)

	//begin actual test

	type document struct {
		DocumentImplementation
		Name string
	}

	var doc document
	doc.Name = "test-document"

	err = docEnd.PostDocument(&doc, docs.Name, nil)

	if err != nil {
		t.Fatal("Unexpected error when creating new document: ", err)
	}

	if doc.Id() == "" {
		t.Fatal("Expected the Id of the document to be set.")
	}

	if doc.Key() == "" {
		t.Fatal("Expected the Key of the document to be set.")
	}

	if doc.Rev() == "" {
		t.Fatal("Expected the Rev of the document to be set.")
	}

	var fetcher document
	err = docEnd.GetDocument(doc.Id(), &fetcher, nil)

	if err != nil {
		t.Fatal("Problems fetching document.")
	}

	if fetcher.Name != doc.Name {
		t.Fatal("Did not fetch Name attribute properly.")
	}

	if fetcher.Id() != doc.Id() {
		t.Fatal("Ids are different between saved and fetched document for no reason.")
	}
}
