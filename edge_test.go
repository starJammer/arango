package arango

import (
	"net/http"
	"testing"
)

func TestMeetsEdgeEndpoint(t *testing.T) {
	var _ EdgeEndpoint = getDatabase("_sysem").EdgeEndpoint()
}

func TestGetEdgesEmptyCollection(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	docs, err := de.GetEdges(opts.Name, nil)

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
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc document
	doc.Name = "test-document"

	err := de.PostEdge(nil, opts.Name, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected to receive error with nil document.",
	)
}

func TestPostEmptyDoc(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc struct {
		EdgeImplementation
	}

	err := de.PostEdge(&doc, opts.Name, nil)

	if err != nil {
		t.Fatal("Unexpected error when posting empty doc: ", err)
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

}

func TestDeleteEdgeBlankName(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.DeleteEdge("", nil)
	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when DeleteEdge blank handler.",
	)
}

func TestDeleteEdgeNonExistent(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.DeleteEdge("non/1", nil)
	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when DeleteEdge non-exstent document.",
	)
}

func TestDeleteEdgeBadName(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.DeleteEdge("bad", nil)
	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when DeleteEdge bad name.",
	)
}

func TestDeleteEdge(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc struct {
		EdgeImplementation
	}

	de.PostEdge(&doc, opts.Name, nil)

	err := de.DeleteEdge(doc.Id(), &DeleteEdgeOptions{Rev: "1"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected 412 error with bad revision.",
	)

	err = de.DeleteEdge(doc.Id(), &DeleteEdgeOptions{IfMatch: "1"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected 412 error with bad IfMatch",
	)

	err = de.DeleteEdge(doc.Id(), nil)

	if err != nil {
		t.Fatal("Unexpected error when deleting document.")
	}

	de.PostEdge(&doc, opts.Name, nil)
	err = de.DeleteEdge(doc.Id(), &DeleteEdgeOptions{Rev: "1", Policy: "last"})
	if err != nil {
		t.Fatal("Unexpected error when deleting with policy = last: ", err)
	}
}

func TestGetEdgesAfterPost(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc struct {
		EdgeImplementation
	}

	de.PostEdge(&doc, opts.Name, nil)

	documents, err := de.GetEdges(
		"test",
		&GetEdgesOptions{ReturnType: "id"},
	)

	if err != nil {
		t.Fatal("Unexpected error when fetching all documents in collection \"test\".")
	}

	if len(documents) != 1 {
		t.Fatal("Expected only one document in collection \"test\": ", documents)
	}

	if documents[0] != doc.Id() {
		t.Fatal("Could not fetch the ids properly. Actual(%s), Expected(%s)")
	}
}

func TestGetEdgeBlankName(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.GetEdge("", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when getting document with blank handle.",
	)
}

func TestGetEdgeNonExistent(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc document
	doc.Name = "test-document"

	de.PostEdge(&doc, opts.Name, nil)
	err := de.GetEdge("non/1234", nil, nil)

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when getting document with blank handle.",
	)
}

func TestGetEdgeBadName(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc document
	doc.Name = "test-document"

	de.PostEdge(&doc, opts.Name, nil)
	err := de.GetEdge("bad", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when getting document with bad handle format.",
	)
}

func TestGetEdgeAfterPost(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc document
	doc.Name = "name"

	de.PostEdge(&doc, opts.Name, nil)
	defer de.DeleteEdge(doc.Id(), nil)

	var fetcher document

	err := de.GetEdge(
		doc.Id(),
		&fetcher,
		nil,
	)

	if err != nil {
		t.Fatal("Unexpected error when fetching a posted document: ", err)
	}

	if fetcher.Id() != doc.Id() {
		t.Fatalf(
			"Fetched document has wrong id: Actual(%s), Expected(%s)",
			fetcher.Id(),
			doc.Id(),
		)
	}

	if fetcher.Key() != doc.Key() {
		t.Fatalf(
			"Fetched document has wrong key: Actual(%s), Expected(%s)",
			fetcher.Key(),
			doc.Key(),
		)
	}

	if fetcher.Rev() != doc.Rev() {
		t.Fatalf(
			"Fetched document has wrong rev: Actual(%s), Expected(%s)",
			fetcher.Rev(),
			doc.Rev(),
		)
	}

	if fetcher.Name != doc.Name {
		t.Fatalf(
			"Fetched document has wrong Name: Actual(%s), Expected(%s)",
			fetcher.Name,
			doc.Name,
		)
	}
}

func TestHeadForBlankName(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	rev, err := de.HeadEdge("", nil)

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

func TestHeadForNonExistent(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	rev, err := de.HeadEdge("none/123434", nil)

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected to receive error for Head with non-existent doc",
	)

	if rev != "" {
		t.Fatal("Expected rev to be blank.")
	}
}

func TestHeadForBadName(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	rev, err := de.HeadEdge("bad", nil)

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

func TestHeadAfterPost(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc document
	doc.Name = "name"

	de.PostEdge(&doc, opts.Name, nil)
	defer de.DeleteEdge(doc.Id(), nil)

	rev, err := de.HeadEdge(doc.Id(), nil)
	if err != nil {
		t.Fatal("Uexpected error")
	}
	if rev != doc.Rev() {
		t.Fatal("Expected rev to equal doc's rev.")
	}

	rev, err = de.HeadEdge(doc.Id(), &HeadEdgeOptions{Rev: doc.Rev()})
	if err != nil {
		t.Fatal("Uexpected error")
	}
	if rev != doc.Rev() {
		t.Fatal("Expected rev to equal doc's rev.")
	}

	rev, err = de.HeadEdge(doc.Id(), &HeadEdgeOptions{IfMatch: doc.Rev()})
	if err != nil {
		t.Fatal("Uexpected error")
	}
	if rev != doc.Rev() {
		t.Fatal("Expected rev to equal doc's rev.")
	}

	rev, err = de.HeadEdge(doc.Id(), &HeadEdgeOptions{IfMatch: "12341234"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected a 412 error with the revision.",
	)
	if rev != doc.Rev() {
		t.Fatal("Expected rev to equal doc's rev.")
	}

	rev, err = de.HeadEdge(doc.Id(), &HeadEdgeOptions{IfNoneMatch: doc.Rev()})
	verifyError(
		err,
		t,
		http.StatusNotModified,
		"Expected a 304 with the revision.",
	)
	if rev != doc.Rev() {
		t.Fatal("Expected rev to equal doc's rev.")
	}

	rev, err = de.HeadEdge(doc.Id(), &HeadEdgeOptions{IfNoneMatch: "12341234123412341234"})
	if err != nil {
		t.Fatal("Uexpected error")
	}
	if rev != doc.Rev() {
		t.Fatal("Expected rev to equal doc's rev.")
	}
}

func TestPutEdgeBlankName(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.PutEdge("", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PutEdge with blank handler.",
	)
}

func TestPutEdgeNonExistent(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.PutEdge("non/1234", &struct{}{}, nil)

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error during PutEdge with non-existent handler.",
	)
}

func TestPutEdgeBadHandler(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.PutEdge("bad", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PutEdge with bad handler.",
	)
}

func TestPutEdgeNilEdge(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc *document = new(document)
	doc.Name = "test"

	de.PostEdge(&doc, opts.Name, nil)
	defer de.DeleteEdge(doc.Id(), nil)

	err := de.PutEdge(doc.Id(), nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PutEdge with blank handler.",
	)
}

func TestPutEdge(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc *document = new(document)
	doc.Name = "test"

	de.PostEdge(&doc, opts.Name, nil)
	defer de.DeleteEdge(doc.Id(), nil)

	var other *document = new(document)
	other.Address = "other"

	err := de.PutEdge(doc.Id(), other, &PutEdgeOptions{Rev: "1"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected error if Rev doesn't match.",
	)

	err = de.PutEdge(doc.Id(), other, &PutEdgeOptions{IfMatch: "1"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected error if IfMatch doesn't match.",
	)

	err = de.PutEdge(doc.Id(), other, nil)
	if err != nil {
		t.Fatal("Unexpected error when putting: ", err)
	}

	var fetcher *document = new(document)
	de.GetEdge(other.Id(), fetcher, nil)

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

	err = de.PutEdge(other.Id(), doc, &PutEdgeOptions{Rev: "12341234", Policy: "last"})
	if err != nil {
		t.Fatal("Unexpected error when putting with policy = last: ", err)
	}
}

func TestPatchEdgeBlankName(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.PatchEdge("", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PatchEdge with blank handler.",
	)
}

func TestPatchEdgeNonExistent(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.PatchEdge("non/1234", &struct{}{}, nil)

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error during PatchEdge with non-existent handler.",
	)
}

func TestPatchEdgeBadHandler(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.PatchEdge("bad", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PatchEdge with bad handler.",
	)
}

func TestPatchEdgeNilEdge(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc *document = new(document)
	doc.Name = "test"

	de.PostEdge(&doc, opts.Name, nil)
	defer de.DeleteEdge(doc.Id(), nil)

	err := de.PatchEdge(doc.Id(), nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PatchEdge with blank handler.",
	)
}

func TestPatchEdge(t *testing.T) {
	var ce = getCE("_system")
	var de = getEE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc *document = new(document)
	doc.Name = "test"

	de.PostEdge(&doc, opts.Name, nil)
	defer de.DeleteEdge(doc.Id(), nil)

	var other *document = new(document)
	other.Address = "other"

	err := de.PatchEdge(doc.Id(), other, &PatchEdgeOptions{Rev: "1111111"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected error if Rev doesn't match.",
	)

	err = de.PatchEdge(doc.Id(), other, &PatchEdgeOptions{IfMatch: "1111111"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected error if IfMatch doesn't match.",
	)

	err = de.PatchEdge(doc.Id(), other, nil)
	if err != nil {
		t.Fatal("Unexpected error when putting: ", err)
	}

	var fetcher *document = new(document)
	de.GetEdge(other.Id(), fetcher, nil)

	if fetcher.Name != doc.Name {
		t.Fatal("Patch failed to preserve  name.")
	}

	if fetcher.Address != other.Address {
		t.Fatalf(
			"Patch failed to set address: Actual(%s), Expected(%s)",
			fetcher.Address,
			other.Address,
		)
	}

	doc.Name = "secondpatch"
	doc.Address = "secondpatch"
	err = de.PatchEdge(other.Id(), doc, &PatchEdgeOptions{Rev: "12341234", Policy: "last"})
	if err != nil {
		t.Fatal("Unexpected error when putting with policy = last: ", err)
	}

	de.GetEdge(doc.Id(), fetcher, nil)

	if fetcher.Name != doc.Name {
		t.Fatal("Patch failed to preserve  name.")
	}

	if fetcher.Address != doc.Address {
		t.Fatalf(
			"Patch failed to set address: Actual(%s), Expected(%s)",
			fetcher.Address,
			doc.Address,
		)
	}
}
