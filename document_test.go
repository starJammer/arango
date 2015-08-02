package arango

import (
	"net/url"
	"testing"
)

func TestPostGetHeadPutPatchDocument(t *testing.T) {
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
		Name    string
		Address string
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

	//attempt to look for a document that shouldn't exist
	var fetcher *document
	fetcher = new(document)
	err = docEnd.GetDocument(doc.Id()+"9999", fetcher, nil)

	if err == nil {
		t.Fatal("Expected an error because the document didn't exist.")
	}

	err = docEnd.GetDocument(doc.Id(), fetcher, nil)

	if err != nil {
		t.Fatal("Problems fetching document.")
	}

	if fetcher.Name != doc.Name {
		t.Fatal("Did not fetch Name attribute properly.")
	}

	if fetcher.Id() != doc.Id() {
		t.Fatal("Ids are different between saved and fetched document for no reason.")
	}

	if fetcher.Rev() != doc.Rev() {
		t.Fatal("Revs are different between saved and fetched document for no reason.")
	}

	err = docEnd.PostDocument(&doc, docs.Name, nil)

	if err == nil {
		t.Fatal("Expected an error when creating a duplicate document but didn't get one.")
	}

	if doc.Id() == "" {
		t.Fatal("Expected the Id of the document to be set even though we double posted.")
	}

	if doc.Key() == "" {
		t.Fatal("Expected the Key of the document to be set even though we double posted.")
	}

	if doc.Rev() == "" {
		t.Fatal("Expected the Rev of the document to be set even though we double posted.")
	}

	err = docEnd.GetDocument(
		doc.Id(),
		fetcher,
		&GetDocumentOptions{IfNoneMatch: doc.Rev()})

	if err != nil {
		t.Fatal("Did not expect an error because the revisions should match.")
	}

	revision, err := docEnd.HeadDocument(doc.Id(), nil)

	if err != nil {
		t.Fatal("Did not expect an error: ", err)
	}

	if revision != doc.Rev() {
		t.Fatalf("Expected HEAD to return correct revision: Expected(%s), Actual(%s)", doc.Rev(), revision)
	}

	revision, err = docEnd.HeadDocument(doc.Id(), &HeadDocumentOptions{IfMatch: doc.Rev() + doc.Rev()})

	if err != nil {
		t.Fatal("Did not expect an error even though revisions did not match:", err, revision)
	}

	if revision != doc.Rev() {
		t.Fatal("Expected returned revision to match document revision.")
	}

	revision, err = docEnd.HeadDocument(doc.Id(), &HeadDocumentOptions{IfNoneMatch: doc.Rev()})

	if err != nil {
		t.Fatal("Did not expect an error: ", err)
	}

	if revision != doc.Rev() {
		t.Fatalf("Expected HEAD to return correct revision: Expected(%s), Actual(%s)", doc.Rev(), revision)
	}

	var newDoc document
	newDoc.Address = "address"

	//test IfMatch
	err = docEnd.PutDocument(doc.Id(), &newDoc, &PutDocumentOptions{IfMatch: doc.Rev() + doc.Rev()})

	if err == nil {
		t.Fatal("Expected an error when putting to a document whose revision doesn't match.")
	}

	//test Rev
	err = docEnd.PutDocument(doc.Id(), &newDoc, &PutDocumentOptions{Rev: doc.Rev() + doc.Rev()})

	if err == nil {
		t.Fatal("Expected an error when putting to a document whose revision doesn't match.")
	}

	err = docEnd.PutDocument(doc.Id(), &newDoc, nil)

	if err != nil {
		t.Fatal("Unexpected error when putting a new document:", err)
	}

	if newDoc.Id() != doc.Id() {
		t.Fatalf("Expected put document id to be equal to old doc id. Expected(%s), Actual(%s)", doc.Id(), newDoc.Id())
	}

	if newDoc.Rev() == doc.Rev() {
		t.Fatalf("Expected put document rev to NOT be  equal to old doc rev. Expected(%s), Actual(%s)", doc.Rev(), newDoc.Rev())
	}

	fetcher = new(document)
	err = docEnd.GetDocument(doc.Id(), fetcher, nil)

	if fetcher.Name != "" || fetcher.Name == "test-document" {
		t.Fatal("Expected fetcher.Name to not have a value since we put a document with no Name. Name = ", fetcher.Name)
	}

	if fetcher.Address != "address" {
		t.Fatalf("Unexpected value for address after put. Actual(%s)", fetcher.Address)
	}

	newDoc.Name = "test-document"
	err = docEnd.PutDocument(doc.Id(), &newDoc, &PutDocumentOptions{Rev: newDoc.Rev() + newDoc.Rev(), Policy: "last"})

	if err != nil {
		t.Fatal("Unexpected error when using \"last\" policy and a mismatching revision: ", err)
	}
}
