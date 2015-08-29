package arango

import (
	"net/http"
	"testing"
)

func TestEEHasDatabase(t *testing.T) {
	var db = getDatabase("_system")
	var ee = db.EdgeEndpoint()

	if ee.Database() == nil {
		t.Fatal("Expected ee to have database reference.")
	}

	if ee.Database().Name() != db.Name() {
		t.Fatal(
			"EE database name(%s), Expected database name (%s)",
			ee.Database().Name(),
			db.Name(),
		)
	}
}

func TestGetEdgesEmptyCollection(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION
	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	docs, err := ee.GetEdges(opts.Name, nil)

	if err != nil {
		t.Fatal("Unexpected err when getting documents: ", err)
	}

	if len(docs) > 0 {
		t.Fatal("Expected no documents in a new collection.")
	}
}

type edge struct {
	EdgeImplementation
	Name    string `json:"Name,omitempty"`
	Address string `json:"Address,omitempty"`
}

func TestPostNilEdge(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	var edge document
	edge.Name = "test-document"

	err := ee.PostEdge(nil, "test/1", "test/1", opts.Name, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected to receive error with nil document.",
	)
}

func TestPostBadFromCollection(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	var edge struct {
		EdgeImplementation
	}

	err := ee.PostEdge(&edge, opts.Name, "fake/1", "test/1", nil)

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error with fake collection name for _from.",
	)
}

func TestPostBadToCollection(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	var edge struct {
		EdgeImplementation
	}

	err := ee.PostEdge(&edge, opts.Name, "test/1", "fake/1", nil)

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error with fake collection name for _to.",
	)
}

func TestPostBlankFrom(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	var edge struct {
		EdgeImplementation
	}

	err := ee.PostEdge(&edge, opts.Name, "", "test/1", nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error with blank _from",
	)
}

func TestPostEmptyEdge(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	var edge struct {
		EdgeImplementation
	}

	err := ee.PostEdge(&edge, opts.Name, "test/1", "test/1", nil)

	if err != nil {
		t.Fatal("Unexpected error when posting empty edge: ", err)
	}

	if edge.Id() == "" {
		t.Fatal("Expected the Id of the document to be set.")
	}

	if edge.Key() == "" {
		t.Fatal("Expected the Key of the document to be set.")
	}

	if edge.Rev() == "" {
		t.Fatal("Expected the Rev of the document to be set.")
	}

	if edge.From() != "test/1" {
		t.Fatalf(
			"Expected _from to be set after a post: Actual(%s), Expected(%s)",
			edge.From(),
			"test/1",
		)
	}

	if edge.To() != "test/1" {
		t.Fatalf(
			"Expected _to to be set after a post: Actual(%s), Expected(%s)",
			edge.To(),
			"test/1",
		)
	}
}

func TestDeleteEdgeBlankName(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	err := ee.DeleteEdge("", nil)
	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when DeleteEdge blank handler.",
	)
}

func TestDeleteEdgeNonExistent(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	err := ee.DeleteEdge("non/1", nil)
	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when DeleteEdge non-exstent document.",
	)
}

func TestDeleteEdgeBadName(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	err := ee.DeleteEdge("bad", nil)
	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when DeleteEdge bad name.",
	)
}

func TestDeleteEdge(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	var edge struct {
		EdgeImplementation
	}

	ee.PostEdge(&edge, opts.Name, "test/1", "test/1", nil)

	err := ee.DeleteEdge(edge.Id(), &DeleteEdgeOptions{Rev: "1"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected 412 error with bad revision.",
	)

	err = ee.DeleteEdge(edge.Id(), &DeleteEdgeOptions{IfMatch: "1"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected 412 error with bad IfMatch",
	)

	err = ee.DeleteEdge(edge.Id(), nil)

	if err != nil {
		t.Fatal("Unexpected error when deleting document.")
	}

	ee.PostEdge(&edge, opts.Name, "test/1", "test/1", nil)
	err = ee.DeleteEdge(edge.Id(), &DeleteEdgeOptions{Rev: "1", Policy: "last"})
	if err != nil {
		t.Fatal("Unexpected error when deleting with policy = last: ", err)
	}
}

func TestGetEdgesAfterPost(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	var edge struct {
		EdgeImplementation
	}

	ee.PostEdge(&edge, opts.Name, "test/1", "test/1", nil)

	documents, err := ee.GetEdges(
		"test",
		&GetEdgesOptions{ReturnType: "id"},
	)

	if err != nil {
		t.Fatal("Unexpected error when fetching all documents in collection \"test\".")
	}

	if len(documents) != 1 {
		t.Fatal("Expected only one document in collection \"test\": ", documents)
	}

	if documents[0] != edge.Id() {
		t.Fatal("Could not fetch the ids properly. Actual(%s), Expected(%s)")
	}
}

func TestGetEdgeBlankName(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	err := ee.GetEdge("", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when getting document with blank handle.",
	)
}

func TestGetEdgeNonExistent(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	var edge document
	edge.Name = "test-document"

	ee.PostEdge(&edge, opts.Name, "test/1", "test/1", nil)
	err := ee.GetEdge("non/1234", nil, nil)

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when getting document with blank handle.",
	)
}

func TestGetEdgeBadName(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	var edge document
	edge.Name = "test-document"

	ee.PostEdge(&edge, opts.Name, "test/1", "test/1", nil)
	err := ee.GetEdge("bad", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when getting document with bad handle format.",
	)
}

func TestGetEdgeAfterPost(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	var edge document
	edge.Name = "name"

	ee.PostEdge(&edge, opts.Name, "test/1", "test/1", nil)
	defer ee.DeleteEdge(edge.Id(), nil)

	var fetcher document

	err := ee.GetEdge(
		edge.Id(),
		&fetcher,
		nil,
	)

	if err != nil {
		t.Fatal("Unexpected error when fetching a posted document: ", err)
	}

	if fetcher.Id() != edge.Id() {
		t.Fatalf(
			"Fetched document has wrong id: Actual(%s), Expected(%s)",
			fetcher.Id(),
			edge.Id(),
		)
	}

	if fetcher.Key() != edge.Key() {
		t.Fatalf(
			"Fetched document has wrong key: Actual(%s), Expected(%s)",
			fetcher.Key(),
			edge.Key(),
		)
	}

	if fetcher.Rev() != edge.Rev() {
		t.Fatalf(
			"Fetched document has wrong rev: Actual(%s), Expected(%s)",
			fetcher.Rev(),
			edge.Rev(),
		)
	}

	if fetcher.Name != edge.Name {
		t.Fatalf(
			"Fetched document has wrong Name: Actual(%s), Expected(%s)",
			fetcher.Name,
			edge.Name,
		)
	}
}

func TestEdgeHeadForBlankName(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	rev, err := ee.HeadEdge("", nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected to receive error for Head with blank name.",
	)

	if rev != "" {
		t.Fatal("Expected rev to be blank.")
	}
}

func TestEdgeHeadForNonExistent(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	rev, err := ee.HeadEdge("none/123434", nil)

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected to receive error for Head with non-existent edge",
	)

	if rev != "" {
		t.Fatal("Expected rev to be blank.")
	}
}

func TestEdgeHeadForBadName(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	rev, err := ee.HeadEdge("bad", nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected to receive error for Head with bad name.",
	)

	if rev != "" {
		t.Fatal("Expected rev to be blank.")
	}
}

func TestEdgeHeadAfterPost(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	var edge document
	edge.Name = "name"

	ee.PostEdge(&edge, opts.Name, "test/1", "test/1", nil)
	defer ee.DeleteEdge(edge.Id(), nil)

	rev, err := ee.HeadEdge(edge.Id(), nil)
	if err != nil {
		t.Fatal("Uexpected error")
	}
	if rev != edge.Rev() {
		t.Fatal("Expected rev to equal edge's rev.")
	}

	rev, err = ee.HeadEdge(edge.Id(), &HeadEdgeOptions{Rev: edge.Rev()})
	if err != nil {
		t.Fatal("Uexpected error")
	}
	if rev != edge.Rev() {
		t.Fatal("Expected rev to equal edge's rev.")
	}

	rev, err = ee.HeadEdge(edge.Id(), &HeadEdgeOptions{IfMatch: edge.Rev()})
	if err != nil {
		t.Fatal("Uexpected error")
	}
	if rev != edge.Rev() {
		t.Fatal("Expected rev to equal edge's rev.")
	}

	rev, err = ee.HeadEdge(edge.Id(), &HeadEdgeOptions{IfMatch: "12341234"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected a 412 error with the revision.",
	)
	if rev != edge.Rev() {
		t.Fatal("Expected rev to equal edge's rev.")
	}

	rev, err = ee.HeadEdge(edge.Id(), &HeadEdgeOptions{IfNoneMatch: edge.Rev()})
	verifyError(
		err,
		t,
		http.StatusNotModified,
		"Expected a 304 with the revision.",
	)
	if rev != edge.Rev() {
		t.Fatal("Expected rev to equal edge's rev.")
	}

	rev, err = ee.HeadEdge(edge.Id(), &HeadEdgeOptions{IfNoneMatch: "12341234123412341234"})
	if err != nil {
		t.Fatal("Uexpected error")
	}
	if rev != edge.Rev() {
		t.Fatal("Expected rev to equal edge's rev.")
	}
}

func TestPutEdgeBlankName(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	err := ee.PutEdge("", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PutEdge with blank handler.",
	)
}

func TestPutEdgeNonExistent(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	err := ee.PutEdge("non/1234", &struct{}{}, nil)

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error during PutEdge with non-existent handler.",
	)
}

func TestPutEdgeBadHandler(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	err := ee.PutEdge("bad", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PutEdge with bad handler.",
	)
}

func TestPutEdgeNilEdge(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	var edge *document = new(document)
	edge.Name = "test"

	ee.PostEdge(&edge, opts.Name, "test/1", "test/1", nil)
	defer ee.DeleteEdge(edge.Id(), nil)

	err := ee.PutEdge(edge.Id(), nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PutEdge with blank handler.",
	)
}

func TestPutEdge(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	var edge *document = new(document)
	edge.Name = "test"

	ee.PostEdge(&edge, opts.Name, "test/1", "test/1", nil)
	defer ee.DeleteEdge(edge.Id(), nil)

	var other *document = new(document)
	other.Address = "other"

	err := ee.PutEdge(edge.Id(), other, &PutEdgeOptions{Rev: "1"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected error if Rev doesn't match.",
	)

	err = ee.PutEdge(edge.Id(), other, &PutEdgeOptions{IfMatch: "1"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected error if IfMatch doesn't match.",
	)

	err = ee.PutEdge(edge.Id(), other, nil)
	if err != nil {
		t.Fatal("Unexpected error when putting: ", err)
	}

	if other.Rev() == "" || other.Rev() == edge.Rev() {
		t.Fatalf(
			"Unexpected value for Rev after putting: Actual(%s), Previous(%s)",
			other.Rev(),
			edge.Rev(),
		)
	}

	var fetcher *document = new(document)
	ee.GetEdge(other.Id(), fetcher, nil)

	if fetcher.Name != "" {
		t.Fatal("Put failed to remove name.")
	}

	if fetcher.Address != other.Address {
		t.Fatalf(
			"Put failed to set address: Actual(%s), Expected(%s)",
			fetcher.Address,
			other.Address,
		)
	}

	err = ee.PutEdge(other.Id(), edge, &PutEdgeOptions{Rev: "12341234", Policy: "last"})
	if err != nil {
		t.Fatal("Unexpected error when putting with policy = last: ", err)
	}
}

func TestPatchEdgeBlankName(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	err := ee.PatchEdge("", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PatchEdge with blank handler.",
	)
}

func TestPatchEdgeNonExistent(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	err := ee.PatchEdge("non/1234", &struct{}{}, nil)

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error during PatchEdge with non-existent handler.",
	)
}

func TestPatchEdgeBadHandler(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	err := ee.PatchEdge("bad", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PatchEdge with bad handler.",
	)
}

func TestPatchEdgeNilEdge(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	var edge *document = new(document)
	edge.Name = "test"

	ee.PostEdge(&edge, opts.Name, "test/1", "test/1", nil)
	defer ee.DeleteEdge(edge.Id(), nil)

	err := ee.PatchEdge(edge.Id(), nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PatchEdge with blank handler.",
	)
}

func TestPatchEdge(t *testing.T) {
	var ce = getCE("_system")
	var ee = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.Type = EDGE_COLLECTION

	ce.PostCollection(opts.Name, opts)
	defer ce.Delete(opts.Name)

	var edge *document = new(document)
	edge.Name = "test"

	ee.PostEdge(&edge, opts.Name, "test/1", "test/1", nil)
	defer ee.DeleteEdge(edge.Id(), nil)

	var other *document = new(document)
	other.Address = "other"

	err := ee.PatchEdge(edge.Id(), other, &PatchEdgeOptions{Rev: "1111111"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected error if Rev doesn't match.",
	)

	err = ee.PatchEdge(edge.Id(), other, &PatchEdgeOptions{IfMatch: "1111111"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected error if IfMatch doesn't match.",
	)

	err = ee.PatchEdge(edge.Id(), other, nil)
	if err != nil {
		t.Fatal("Unexpected error when putting: ", err)
	}

	if other.Rev() == "" || other.Rev() == edge.Rev() {
		t.Fatalf(
			"Unexpected value for Rev after patching: Actual(%s), Previous(%s)",
			other.Rev(),
			edge.Rev(),
		)
	}

	var fetcher *document = new(document)
	ee.GetEdge(other.Id(), fetcher, nil)

	if fetcher.Name != edge.Name {
		t.Fatal("Patch failed to preserve  name.")
	}

	if fetcher.Address != other.Address {
		t.Fatalf(
			"Patch failed to set address: Actual(%s), Expected(%s)",
			fetcher.Address,
			other.Address,
		)
	}

	edge.Name = "secondpatch"
	edge.Address = "secondpatch"
	err = ee.PatchEdge(other.Id(), edge, &PatchEdgeOptions{Rev: "12341234", Policy: "last"})
	if err != nil {
		t.Fatal("Unexpected error when putting with policy = last: ", err)
	}

	ee.GetEdge(edge.Id(), fetcher, nil)

	if fetcher.Name != edge.Name {
		t.Fatal("Patch failed to preserve  name.")
	}

	if fetcher.Address != edge.Address {
		t.Fatalf(
			"Patch failed to set address: Actual(%s), Expected(%s)",
			fetcher.Address,
			edge.Address,
		)
	}
}
