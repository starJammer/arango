package arango

import (
	"net/http"
	"testing"
)

func TestByExampleBlankCollection(t *testing.T) {
	var se = getSE("_system")

	cursor, err := se.ByExample("", map[string]int{"count": -1}, nil)

	if cursor != nil {
		t.Fatal("Expected a nil cursor.")
	}

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected 404 with blank collection name.",
	)
}

func TestByExampleGetsNothing(t *testing.T) {
	byExampleDataInit()
	var se = getSE("_system")

	cursor, err := se.ByExample("test", map[string]int{"count": -1}, nil)

	if err != nil {
		t.Fatal("Uexpected error when searching by example: ", err)
	}

	if cursor.HasMore() {
		t.Fatal("Unexpected : cursor.HasMore() == true")
	}

	if cursor.Count() > 0 {
		t.Fatal("Unexpected : cursor.Count() > 0")
	}
}

func TestByExampleFetchesOneOnly(t *testing.T) {
	var se = getSE("_system")

	cursor, err := se.ByExample("test", map[string]int{"count": 1}, nil)

	if err != nil {
		t.Fatal("Uexpected error when searching by example: ", err)
	}

	if !cursor.HasMore() {
		t.Fatal("Expected : cursor.HasMore() == true")
	}

	if cursor.Count() != 1 {
		t.Fatal("Expected : cursor.Count() == 1")
	}

	var fetcher countDocument
	err = cursor.Next(&fetcher)

	if err != nil {
		t.Fatal("Unexpected error when reading Next:", err)
	}

	if fetcher.Count != 1 {
		t.Fatal("Fetched the wrong document.")
	}

	if fetcher.Letter != "a" {
		t.Fatal("Expected letter to be filled.")
	}

	if cursor.HasMore() {
		t.Fatal("Expected : cursor.HasMore() == false")
	}

	err = cursor.Next(&fetcher)

	if err != EOC {
		t.Fatal("Expected end of cursor error when trying to read next.")
	}

}

func TestFinalTest(t *testing.T) {
	destroyByExampleDataInit()
}

type countDocument struct {
	DocumentImplementation
	Count  int    `json:"count"`
	Letter string `json:"letter"`
}

func byExampleDataInit() {
	var ce = getCE("_system")
	var de = getDE("_system")

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	ce.PostCollection(opts.Name, nil)

	for i := 0; i < 10; i++ {
		doc := &countDocument{
			Count:  i,
			Letter: "a",
		}
		if i >= 4 {
			doc.Letter = "b"
		}

		de.PostDocument(
			doc,
			opts.Name,
			nil,
		)
	}
}

func destroyByExampleDataInit() {
	var ce = getCE("_system")
	ce.Delete("test")
}
