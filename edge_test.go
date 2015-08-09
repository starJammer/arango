package arango

import (
	"net/url"
	"testing"
)

func TestPostGetHeadPutPatchEdge(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")
	var collEnd = db.CollectionEndpoint()
	var docEnd = db.DocumentEndpoint()
	var edgeEnd = db.EdgeEndpoint()

	docs := DefaultPostCollectionOptions()
	docs.Name = "test"

	collEnd.PostCollection(docs.Name, nil)
	defer collEnd.Delete(docs.Name)

	edgeCollectionOptions := DefaultPostCollectionOptions()
	edgeCollectionOptions.Type = EDGE_COLLECTION
	edgeCollectionOptions.Name = "edge-test"

	collEnd.PostCollection(edgeCollectionOptions.Name, edgeCollectionOptions)
	defer collEnd.Delete(edgeCollectionOptions.Name)
	//begin actual test
	type document struct {
		DocumentImplementation
		Name string
	}
	type edge struct {
		EdgeImplementation
		Name    string
		Address string
	}

	edges, err := edgeEnd.GetEdges(edgeCollectionOptions.Name, nil)
	if err != nil {
		t.Fatal("Unexpected error when fetching all edges in collection.")
	}
	if len(edges) != 0 {
		t.Fatal("Expected no edges in collection.")
	}

	var doc document
	doc.Name = "test-document"

	err = docEnd.PostDocument(&doc, docs.Name, nil)

	if err != nil {
		t.Fatal("Unexpected error when creating new document: ", err)
	}

	var efetcher *edge
	efetcher = new(edge)

	//create an edge from/to the same document
	var dedge *edge
	dedge = new(edge)
	dedge.Name = "edge-1"
	dedge.Address = "edge-address"
	err = edgeEnd.PostEdge(dedge, edgeCollectionOptions.Name, doc.Id(), doc.Id(), nil)

	if err != nil {
		t.Fatal("Unexpected error when saving the edge.")
	}

	if dedge.Id() == "" {
		t.Fatal("Expected edge id to be set.")
	}

	if dedge.Key() == "" {
		t.Fatal("Expected edge key to be set.")
	}

	if dedge.Rev() == "" {
		t.Fatal("Expected edge revision to be set.")
	}

	if dedge.From() != doc.Id() || dedge.To() != doc.Id() {
		t.Fatalf("Unexpected From and To fields in the edge: From(%s), To(%s), Expected From(%s), Expected To(%s)", dedge.From(), dedge.To(), doc.Id(), doc.Id())
	}

	err = edgeEnd.GetEdge(dedge.Id(), efetcher, nil)

	if err != nil {
		t.Fatal("Unexpected error when fetching the edge again for some reason.")
	}

	if efetcher.From() != dedge.From() || efetcher.To() != dedge.To() {
		t.Fatalf("Unexpected From and To fields in fetched edge. From(%s), To(%s), Expected From(%s), Expected To(%s)", efetcher.From(), efetcher.To(), doc.Id(), doc.Id())
	}

	if efetcher.Name != dedge.Name {
		t.Fatalf("Expected Name attributes to match. Actual(%s), Expected(%s)", efetcher.Name, dedge.Name)
	}

	edges, err = edgeEnd.GetEdges(edgeCollectionOptions.Name, nil)
	if err != nil {
		t.Fatal("Unexpected error when fetching all edges in collection.")
	}
	if len(edges) != 1 {
		t.Fatal("Expected one edge in collection.")
	}

	err = edgeEnd.GetEdge(
		dedge.Id(),
		efetcher,
		&GetEdgeOptions{IfNoneMatch: dedge.Rev()})

	if err != nil {
		t.Fatal("Did not expect an error because the revisions should match.")
	}

	revision, err := edgeEnd.HeadEdge(dedge.Id(), nil)

	if err != nil {
		t.Fatal("Did not expect an error: ", err)
	}

	if revision != dedge.Rev() {
		t.Fatalf("Expected HEAD to return correct revision: Expected(%s), Actual(%s)", doc.Rev(), revision)
	}

	revision, err = edgeEnd.HeadEdge(dedge.Id(), &HeadEdgeOptions{IfMatch: dedge.Rev() + dedge.Rev()})

	if err != nil {
		t.Fatal("Did not expect an error even though revisions did not match:", err, revision)
	}

	if revision != dedge.Rev() {
		t.Fatal("Expected returned revision to match edge revision.")
	}

	revision, err = edgeEnd.HeadEdge(dedge.Id(), &HeadEdgeOptions{IfNoneMatch: dedge.Rev()})

	if err != nil {
		t.Fatal("Did not expect an error: ", err)
	}

	if revision != dedge.Rev() {
		t.Fatalf("Expected HEAD to return correct revision: Expected(%s), Actual(%s)", dedge.Rev(), revision)
	}

	var newEdge edge
	newEdge.Address = "address"

	//test IfMatch
	err = edgeEnd.PutEdge(dedge.Id(), &newEdge, &PutEdgeOptions{IfMatch: dedge.Rev() + dedge.Rev()})

	if err == nil {
		t.Fatal("Expected an error when putting to a document whose revision doesn't match.")
	}

	//test Rev
	err = edgeEnd.PutEdge(dedge.Id(), &newEdge, &PutEdgeOptions{Rev: dedge.Rev() + dedge.Rev()})

	if err == nil {
		t.Fatal("Expected an error when putting to a document whose revision doesn't match.")
	}

	err = edgeEnd.PutEdge(dedge.Id(), &newEdge, nil)

	if err != nil {
		t.Fatal("Unexpected error when putting a new document:", err)
	}

	if newEdge.Id() != dedge.Id() {
		t.Fatalf("Expected put document id to be equal to old doc id. Expected(%s), Actual(%s)", dedge.Id(), newEdge.Id())
	}

	if newEdge.Rev() == dedge.Rev() {
		t.Fatalf("Expected put document rev to NOT be  equal to old doc rev. Expected(%s), Actual(%s)", dedge.Rev(), newEdge.Rev())
	}

	if newEdge.From() != dedge.From() || newEdge.To() != dedge.To() {
		t.Fatalf("Expected put edge to have correct from and to.. ExpectedFrom(%s), ExpectedTo(%s), ActualFrom(%s), ActualTo(%s)", newEdge.From(), newEdge.To(), dedge.From(), dedge.To())
	}

	efetcher = new(edge)
	err = edgeEnd.GetEdge(dedge.Id(), efetcher, nil)

	if efetcher.Name != "" || efetcher.Name == "test-document" {
		t.Fatal("Expected efetcher.Name to not have a value since we put a document with no Name. Name = ", efetcher.Name)
	}

	if efetcher.Address != "address" {
		t.Fatalf("Unexpected value for address after put. Actual(%s)", efetcher.Address)
	}

	newEdge.Name = "test-document"
	err = edgeEnd.PutEdge(dedge.Id(), &newEdge, &PutEdgeOptions{Rev: newEdge.Rev() + newEdge.Rev(), Policy: "last"})

	if err != nil {
		t.Fatal("Unexpected error when using \"last\" policy and a mismatching revision: ", err)
	}

	efetcher = new(edge)
	err = edgeEnd.GetEdge(dedge.Id(), efetcher, nil)

	if efetcher.Name != newEdge.Name || efetcher.Address != newEdge.Address {
		t.Fatalf("Got unexpected values after putting a new document using last policy: Expected(%v), Actual(%v)", newEdge, efetcher)
	}

	var patcher struct {
		Name string
	}

	patcher.Name = "new-name"

	err = edgeEnd.PatchEdge(dedge.Id(), &patcher, nil)

	if err != nil {
		t.Fatal("Unexpected error when patching using a map: ", err)
	}

	efetcher = new(edge)
	err = edgeEnd.GetEdge(dedge.Id(), efetcher, nil)

	if efetcher.Name != "new-name" || efetcher.Address != newEdge.Address {
		t.Fatalf("Unexpected error when patching only one field. Expected-field-value(%v), Actual(%v)", newEdge.Name, efetcher.Name)
	}

	//test patching with a map
	err = edgeEnd.PatchEdge(dedge.Id(), &map[string]interface{}{"Name": "new-name"}, nil)

	if err != nil {
		t.Fatal("Unexpected error when patching using a map: ", err)
	}

	efetcher = new(edge)
	err = edgeEnd.GetEdge(dedge.Id(), efetcher, nil)

	if efetcher.Name != "new-name" || efetcher.Address != newEdge.Address {
		t.Fatalf("Unexpected error when patching with a map. Expected-field-value(%v), Actual(%v)", newEdge.Name, efetcher.Name)
	}

	err = edgeEnd.DeleteEdge(dedge.Id(), &DeleteEdgeOptions{IfMatch: efetcher.Rev() + efetcher.Rev()})

	if err == nil {
		t.Fatal("Expected delete with bad revision to fail.")
	}

	err = edgeEnd.DeleteEdge(dedge.Id(), &DeleteEdgeOptions{IfMatch: efetcher.Rev()})

	if err != nil {
		t.Fatal("Unexpected error when deleting document: ", err)
	}
}
