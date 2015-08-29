package arango

import (
	"net/http"
	"testing"
)

func TestDEHasDatabase(t *testing.T) {
	var db = getDatabase("_system")
	var de = db.DocumentEndpoint()

	if de.Database() == nil {
		t.Fatal("Expected ee to have database reference.")
	}

	if de.Database().Name() != db.Name() {
		t.Fatal(
			"DE database name(%s), Expected database name (%s)",
			de.Database().Name(),
			db.Name(),
		)
	}
}

func TestGetDocumentsEmptyCollection(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	docs, err := de.GetDocuments(opts.Name, nil)

	if err != nil {
		t.Fatal("Unexpected err when getting documents: ", err)
	}

	if len(docs) > 0 {
		t.Fatal("Expected no documents in a new collection.")
	}
}

type document struct {
	DocumentImplementation
	Name    string `json:"Name,omitempty"`
	Address string `json:"Address,omitempty"`
}

func TestPostNilDocument(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc document
	doc.Name = "test-document"

	err := de.PostDocument(nil, opts.Name, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected to receive error with nil document.",
	)
}

func TestPostEmptyDoc(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc struct {
		DocumentImplementation
	}

	err := de.PostDocument(&doc, opts.Name, nil)

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

func TestDeleteDocumentBlankName(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.DeleteDocument("", nil)
	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when DeleteDocument blank handler.",
	)
}

func TestDeleteDocumentNonExistent(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.DeleteDocument("non/1", nil)
	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when DeleteDocument non-exstent document.",
	)
}

func TestDeleteDocumentBadName(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.DeleteDocument("bad", nil)
	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when DeleteDocument bad name.",
	)
}

func TestDeleteDocument(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc struct {
		DocumentImplementation
	}

	de.PostDocument(&doc, opts.Name, nil)

	err := de.DeleteDocument(doc.Id(), &DeleteDocumentOptions{Rev: "1"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected 412 error with bad revision.",
	)

	err = de.DeleteDocument(doc.Id(), &DeleteDocumentOptions{IfMatch: "1"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected 412 error with bad IfMatch",
	)

	err = de.DeleteDocument(doc.Id(), nil)

	if err != nil {
		t.Fatal("Unexpected error when deleting document.")
	}

	de.PostDocument(&doc, opts.Name, nil)
	err = de.DeleteDocument(doc.Id(), &DeleteDocumentOptions{Rev: "1", Policy: "last"})
	if err != nil {
		t.Fatal("Unexpected error when deleting with policy = last: ", err)
	}
}

func TestGetDocumentsAfterPost(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc struct {
		DocumentImplementation
	}

	de.PostDocument(&doc, opts.Name, nil)

	documents, err := de.GetDocuments(
		"test",
		&GetDocumentsOptions{ReturnType: "id"},
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

func TestGetDocumentBlankName(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.GetDocument("", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when getting document with blank handle.",
	)
}

func TestGetDocumentNonExistent(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc document
	doc.Name = "test-document"

	de.PostDocument(&doc, opts.Name, nil)
	err := de.GetDocument("non/1234", nil, nil)

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when getting document with blank handle.",
	)
}

func TestGetDocumentBadName(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc document
	doc.Name = "test-document"

	de.PostDocument(&doc, opts.Name, nil)
	err := de.GetDocument("bad", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when getting document with bad handle format.",
	)
}

func TestGetDocumentAfterPost(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc document
	doc.Name = "name"

	de.PostDocument(&doc, opts.Name, nil)
	defer de.DeleteDocument(doc.Id(), nil)

	var fetcher document

	err := de.GetDocument(
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
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	rev, err := de.HeadDocument("", nil)

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
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	rev, err := de.HeadDocument("none/123434", nil)

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
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	rev, err := de.HeadDocument("bad", nil)

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
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc document
	doc.Name = "name"

	de.PostDocument(&doc, opts.Name, nil)
	defer de.DeleteDocument(doc.Id(), nil)

	rev, err := de.HeadDocument(doc.Id(), nil)
	if err != nil {
		t.Fatal("Uexpected error")
	}
	if rev != doc.Rev() {
		t.Fatal("Expected rev to equal doc's rev.")
	}

	rev, err = de.HeadDocument(doc.Id(), &HeadDocumentOptions{Rev: doc.Rev()})
	if err != nil {
		t.Fatal("Uexpected error")
	}
	if rev != doc.Rev() {
		t.Fatal("Expected rev to equal doc's rev.")
	}

	rev, err = de.HeadDocument(doc.Id(), &HeadDocumentOptions{IfMatch: doc.Rev()})
	if err != nil {
		t.Fatal("Uexpected error")
	}
	if rev != doc.Rev() {
		t.Fatal("Expected rev to equal doc's rev.")
	}

	rev, err = de.HeadDocument(doc.Id(), &HeadDocumentOptions{IfMatch: "12341234"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected a 412 error with the revision.",
	)
	if rev != doc.Rev() {
		t.Fatal("Expected rev to equal doc's rev.")
	}

	rev, err = de.HeadDocument(doc.Id(), &HeadDocumentOptions{IfNoneMatch: doc.Rev()})
	verifyError(
		err,
		t,
		http.StatusNotModified,
		"Expected a 304 with the revision.",
	)
	if rev != doc.Rev() {
		t.Fatal("Expected rev to equal doc's rev.")
	}

	rev, err = de.HeadDocument(doc.Id(), &HeadDocumentOptions{IfNoneMatch: "12341234123412341234"})
	if err != nil {
		t.Fatal("Uexpected error")
	}
	if rev != doc.Rev() {
		t.Fatal("Expected rev to equal doc's rev.")
	}
}

func TestPutDocumentBlankName(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.PutDocument("", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PutDocument with blank handler.",
	)
}

func TestPutDocumentNonExistent(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.PutDocument("non/1234", &struct{}{}, nil)

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error during PutDocument with non-existent handler.",
	)
}

func TestPutDocumentBadHandler(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.PutDocument("bad", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PutDocument with bad handler.",
	)
}

func TestPutDocumentNilDocument(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc *document = new(document)
	doc.Name = "test"

	de.PostDocument(&doc, opts.Name, nil)
	defer de.DeleteDocument(doc.Id(), nil)

	err := de.PutDocument(doc.Id(), nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PutDocument with blank handler.",
	)
}

func TestPutDocument(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc *document = new(document)
	doc.Name = "test"

	de.PostDocument(&doc, opts.Name, nil)
	defer de.DeleteDocument(doc.Id(), nil)

	var other *document = new(document)
	other.Address = "other"

	err := de.PutDocument(doc.Id(), other, &PutDocumentOptions{Rev: "1"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected error if Rev doesn't match.",
	)

	err = de.PutDocument(doc.Id(), other, &PutDocumentOptions{IfMatch: "1"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected error if IfMatch doesn't match.",
	)

	err = de.PutDocument(doc.Id(), other, nil)
	if err != nil {
		t.Fatal("Unexpected error when putting: ", err)
	}

	if other.Rev() == "" || other.Rev() == doc.Rev() {
		t.Fatalf(
			"Unexpected value for Rev after putting: Actual(%s), Previous(%s)",
			other.Rev(),
			doc.Rev(),
		)
	}

	var fetcher *document = new(document)
	de.GetDocument(other.Id(), fetcher, nil)

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

	err = de.PutDocument(other.Id(), doc, &PutDocumentOptions{Rev: "12341234", Policy: "last"})
	if err != nil {
		t.Fatal("Unexpected error when putting with policy = last: ", err)
	}
}

func TestPatchDocumentBlankName(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.PatchDocument("", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PatchDocument with blank handler.",
	)
}

func TestPatchDocumentNonExistent(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.PatchDocument("non/1234", &struct{}{}, nil)

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error during PatchDocument with non-existent handler.",
	)
}

func TestPatchDocumentBadHandler(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	err := de.PatchDocument("bad", nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PatchDocument with bad handler.",
	)
}

func TestPatchDocumentNilDocument(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc *document = new(document)
	doc.Name = "test"

	de.PostDocument(&doc, opts.Name, nil)
	defer de.DeleteDocument(doc.Id(), nil)

	err := de.PatchDocument(doc.Id(), nil, nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error during PatchDocument with blank handler.",
	)
}

func TestPatchDocument(t *testing.T) {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)
	defer ce.Delete(opts.Name)

	var doc *document = new(document)
	doc.Name = "test"

	de.PostDocument(&doc, opts.Name, nil)
	defer de.DeleteDocument(doc.Id(), nil)

	var other *document = new(document)
	other.Address = "other"

	err := de.PatchDocument(doc.Id(), other, &PatchDocumentOptions{Rev: "1111111"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected error if Rev doesn't match.",
	)

	err = de.PatchDocument(doc.Id(), other, &PatchDocumentOptions{IfMatch: "1111111"})
	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected error if IfMatch doesn't match.",
	)

	err = de.PatchDocument(doc.Id(), other, nil)
	if err != nil {
		t.Fatal("Unexpected error when putting: ", err)
	}

	if other.Rev() == "" || other.Rev() == doc.Rev() {
		t.Fatalf(
			"Unexpected value for Rev after patching: Actual(%s), Previous(%s)",
			other.Rev(),
			doc.Rev(),
		)
	}

	var fetcher *document = new(document)
	de.GetDocument(other.Id(), fetcher, nil)

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
	err = de.PatchDocument(other.Id(), doc, &PatchDocumentOptions{Rev: "12341234", Policy: "last"})
	if err != nil {
		t.Fatal("Unexpected error when putting with policy = last: ", err)
	}

	de.GetDocument(doc.Id(), fetcher, nil)

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
