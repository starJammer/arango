package arango

import (
	"net/http"
	"testing"
)

type simpleDoc struct {
	EdgeImplementation
	Text   string `json:"text"`
	Number int    `json:"number"`
}

func TestDEHasDatabase(t *testing.T) {
	var db = getDatabase("_system")
	var de = db.DocumentEndpoint()

	if de.Database() == nil {
		t.Fatal("Expected document endpoint to have database reference.")
	}

	if de.Database().Name() != db.Name() {
		t.Fatal(
			"DE database name(%s), Expected database name (%s)",
			de.Database().Name(),
			db.Name(),
		)
	}
}

func TestPostWithNilParams(t *testing.T) {
	ce, testName := createTestCollection()
	defer ce.Delete(testName)

	de := getDE("_system")

	err := de.PostDocuments(nil)
	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when posting with nil params",
	)
}

func TestPostEmptyBody(t *testing.T) {
	ce, testName := createTestCollection()
	defer ce.Delete(testName)

	de := getDE("_system")

	postOpts := &PostDocumentOptions{
		Collection: testName,
	}

	err := de.PostDocuments(postOpts)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected to receive error when posting an empty document",
	)
}

func TestPostEmptyCollectionName(t *testing.T) {
	de := getDE("_system")
	postOpts := &PostDocumentOptions{
		Document: &simpleDoc{
			Text: "hi",
		},
	}

	err := de.PostDocuments(postOpts)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected to receive error because of blank collection name",
	)
}

func TestPostBadCollectionName(t *testing.T) {
	de := getDE("_system")
	postOpts := &PostDocumentOptions{
		Collection: "noexist",
		Document: &simpleDoc{
			Text: "hi",
		},
	}

	err := de.PostDocuments(postOpts)

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected to receive error because bad collection name is not found",
	)
}

func TestPostOneDocument(t *testing.T) {
	ce, testName := createTestCollection()
	defer ce.Delete(testName)

	de := getDE("_system")

	var doc = simpleDoc{
		Text: "hi",
	}

	postOpts := &PostDocumentOptions{
		Document:   &doc,
		Collection: testName,
	}

	err := de.PostDocuments(postOpts)

	if err != nil {
		t.Fatal("Unexpected err when posting one document: ", err)
	}

	msg := "Exected arango attributes for document to be populated: "
	if doc.ArangoId == "" {
		t.Fatal(msg, doc)
	}

	if doc.ArangoKey == "" {
		t.Fatal(msg, doc)
	}

	if doc.ArangoRev == "" {
		t.Fatal(msg, doc)
	}
}

func TestPostTwiceWithSameKeyFails(t *testing.T) {
	ce, testName := createTestCollection()
	defer ce.Delete(testName)

	de := getDE("_system")

	var doc = simpleDoc{
		Text: "hi",
	}

	postOpts := &PostDocumentOptions{
		Document:   &doc,
		Collection: testName,
	}

	err := de.PostDocuments(postOpts)
	if err != nil {
		t.Fatal(err)
	}
	//now the document has an ArangoKey
	//that exists
	err = de.PostDocuments(postOpts)
	verifyError(
		err,
		t,
		http.StatusConflict,
		"Posting a document with the same _key as an existing doc should generated an error",
	)
}

func TestPostTwiceWithNoKeyGeneratesNewId(t *testing.T) {
	ce, testName := createTestCollection()
	defer ce.Delete(testName)

	de := getDE("_system")

	var doc = simpleDoc{
		Text: "hi",
	}

	postOpts := &PostDocumentOptions{
		Document:   &doc,
		Collection: testName,
	}

	err := de.PostDocuments(postOpts)
	if err != nil {
		t.Fatal(err)
	}
	firstId := doc.ArangoId
	doc.ArangoKey = ""
	err = de.PostDocuments(postOpts)
	if firstId == doc.ArangoId {
		t.Fatal("Expected id from first post to be different from second posting", firstId, doc)
	}
}

func TestPostMultipleDocuments(t *testing.T) {
	ce, testName := createTestCollection()
	defer ce.Delete(testName)

	de := getDE("_system")

	var doc1 = simpleDoc{
		Text: "hi",
	}

	var doc2 = simpleDoc{
		Text: "bye",
	}

	postOpts := &PostDocumentOptions{
		MultiDocuments: []interface{}{&doc1, &doc2},
		Collection:     testName,
	}

	err := de.PostDocuments(postOpts)

	if err != nil {
		t.Fatal("Unexpected err when posting one document: ", err)
	}

	msg := "Exected arango attributes for document to be populated: "

	if doc1.ArangoId == "" || doc2.ArangoId == "" || doc1.ArangoId == doc2.ArangoId {
		t.Fatal(msg, doc1, doc2)
	}

	if doc1.ArangoKey == "" || doc2.ArangoKey == "" || doc1.ArangoKey == doc2.ArangoKey {
		t.Fatal(msg, doc1, doc2)
	}

	if doc1.ArangoRev == "" || doc2.ArangoRev == "" || doc1.ArangoRev == doc2.ArangoRev {
		t.Fatal(msg, doc1, doc2)
	}

	var docs = []interface{}{
		&simpleDoc{
			Text: "1",
		},
		&simpleDoc{
			Text: "2",
		},
	}

	postOpts.MultiDocuments = docs
	err = de.PostDocuments(postOpts)

	if err != nil {
		t.Fatal("Unexpected err when posting one document: ", err)
	}

	for _, v := range docs {
		doc, ok := v.(*simpleDoc)

		if !ok {
			t.Fatal("Unexpected item type", doc)
		}

		if doc.ArangoId == "" {
			t.Fatal(msg, docs)
		}

		if doc.ArangoKey == "" {
			t.Fatal(msg, docs)
		}

		if doc.ArangoRev == "" {
			t.Fatal(msg, docs)
		}
	}
}

func TestGetWithNilParams(t *testing.T) {
	ce, testName := createTestCollection()
	defer ce.Delete(testName)

	de := getDE("_system")

	err := de.GetDocument(nil)

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected an error with nil get params",
	)
}

func TestGetBlankDocumentHandleFails(t *testing.T) {
	ce, testName := createTestCollection()
	defer ce.Delete(testName)

	de := getDE("_system")

	err := de.GetDocument(&GetDocumentOptions{})

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when looking for blank document handle",
	)

}

func TestGetNonExistentDocumentFails(t *testing.T) {
	ce, testName := createTestCollection()
	defer ce.Delete(testName)

	de := getDE("_system")

	err := de.GetDocument(&GetDocumentOptions{
		Handle: "test/1234",
	})

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when looking for a bad document handle",
	)

}

func TestPostDocumentThenGetIt(t *testing.T) {
	ce, testName := createTestCollection()
	defer ce.Delete(testName)

	de := getDE("_system")
	var doc = simpleDoc{
		Text: "text",
	}
	var receiver simpleDoc
	opts := &PostDocumentOptions{
		Collection: testName,
		Document:   &doc,
	}

	err := de.PostDocuments(opts)
	if err != nil {
		t.Fatal(err)
	}

	err = de.GetDocument(&GetDocumentOptions{
		Handle:   doc.ArangoId,
		Document: &receiver,
	})

	if err != nil {
		t.Fatal(err)
	}

	msg := "Expected the original doc and fetched doc to hold the same data"
	if receiver.ArangoId != doc.ArangoId {
		t.Fatal(msg, receiver, doc)
	}

	if receiver.Text != doc.Text {
		t.Fatal(msg, receiver, doc)
	}
}

func TestPostToEdgeCollection(t *testing.T) {
	docC, testName := createTestCollection()
	defer docC.Delete(testName)
	edgeC, testEdge := createTestEdgeCollection()
	defer edgeC.Delete(testEdge)
	de := getDE("_system")

	var doc simpleDoc
	var edge simpleDoc

	de.PostDocuments(&PostDocumentOptions{
		Document:   &doc,
		Collection: testName,
	})

	err := de.PostDocuments(&PostDocumentOptions{
		Document:   &edge,
		Collection: testEdge,
	})

	if err == nil {
		t.Fatal("Expected error when creating an edge with no _from and _to attributes")
	}

	edge.ArangoFrom = doc.ArangoId
	edge.ArangoTo = doc.ArangoId

	err = de.PostDocuments(&PostDocumentOptions{
		Document:   &edge,
		Collection: testEdge,
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteDocument(t *testing.T) {
	ce, testName := createTestCollection()
	defer ce.Delete(testName)

	de := getDE("_system")

	var doc = simpleDoc{
		Text: "test",
	}

	var receiver simpleDoc

	de.PostDocuments(&PostDocumentOptions{
		Document:   &doc,
		Collection: testName,
	})

	err := de.DeleteDocument(&DeleteDocumentOptions{
		Handle:      doc.ArangoId,
		IfMatch:     doc.ArangoRev + "1",
		OldReceiver: &receiver,
	})

	verifyError(
		err,
		t,
		http.StatusPreconditionFailed,
		"Expected error because of bad revision during a delete",
	)

	err = de.DeleteDocument(&DeleteDocumentOptions{
		Handle:      doc.ArangoId,
		IfMatch:     doc.ArangoRev,
		OldReceiver: &receiver,
	})

	if err != nil {
		t.Fatal(err)
	}

	msg := "ReceiveOld failed to retrieve the old document during a delete"
	if receiver.ArangoId != doc.ArangoId {
		t.Fatal(msg, receiver, doc)
	}

	if receiver.Text != doc.Text {
		t.Fatal(msg, receiver, doc)
	}

	err = de.GetDocument(&GetDocumentOptions{
		Handle: doc.ArangoId,
	})

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected the document to be deleted",
	)
}

func TestDeleteMultiDocuments(t *testing.T) {
	ce, testName := createTestCollection()
	defer ce.Delete(testName)

	de := getDE("_system")

	var doc1 = simpleDoc{
		Text: "1",
	}
	var doc2 = simpleDoc{
		Text: "2",
	}

	var doc3 simpleDoc
	var doc4 simpleDoc
	var receiver = []interface{}{
		&doc3,
		&doc4,
	}

	de.PostDocuments(&PostDocumentOptions{
		Document:   &doc1,
		Collection: testName,
	})

	de.PostDocuments(&PostDocumentOptions{
		Document:   &doc2,
		Collection: testName,
	})

	err := de.DeleteMultiDocuments(&DeleteMultiDocumentsOptions{
		Handles: []interface{}{
			doc1.ArangoId,
			doc2.ArangoKey,
		},
		ReturnOld:   true,
		OldReceiver: receiver,
	})

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when deleting multi with no collection name",
	)

	err = de.DeleteMultiDocuments(&DeleteMultiDocumentsOptions{
		Collection: testName,
		Handles: []interface{}{
			doc1.ArangoId,
			doc2.ArangoKey,
		},
		ReturnOld:   true,
		OldReceiver: receiver,
	})

	if err != nil {
		t.Fatal(err)
	}

	msg := "ReturnOld failed to fetch old documents"

	if doc1.ArangoId != doc3.ArangoId {
		t.Fatal(msg, doc1, doc3)
	}
	if doc2.ArangoId != doc4.ArangoId {
		t.Fatal(msg, doc2, doc4)
	}

	if doc1.Text != doc3.Text {
		t.Fatal(msg, doc1, doc3)
	}
	if doc2.Text != doc4.Text {
		t.Fatal(msg, doc2, doc4)
	}
}

func TestPatchDocument(t *testing.T) {
	ce, testName := createTestCollection()
	defer ce.Delete(testName)

	de := getDE("_system")

	var doc1 = simpleDoc{
		Text: "1",
	}

	var oldReceiver simpleDoc
	var newReceiver simpleDoc

	de.PostDocuments(&PostDocumentOptions{
		Document:   &doc1,
		Collection: testName,
	})

	var patcher = simpleDoc{
		Text:   "",
		Number: 10,
	}

	patch := DefaultPatchDocumentOptions()
	patch.Handle = doc1.ArangoId
	patch.Document = &patcher
	//patch.ReturnOld = true
	patch.OldReceiver = &oldReceiver

	//patch.ReturnNew = true
	patch.NewReceiver = &newReceiver

	err := de.PatchDocument(patch)

	if err != nil {
		t.Fatal(err)
	}

	if oldReceiver.Text != "1" {
		t.Fatal("Expected old document to have correct text", oldReceiver)
	}
	if oldReceiver.Number != 0 {
		t.Fatal("Expected old document to have correct number", oldReceiver)
	}
	if newReceiver.Text != "" {
		t.Fatal("Expected new document to have correct text", newReceiver)
	}
	if newReceiver.Number != 10 {
		t.Fatal("Expected new document to have correct number", newReceiver)
	}

}
