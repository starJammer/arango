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

func TestByExampleBadCollection(t *testing.T) {
	var se = getSE("_system")

	cursor, err := se.ByExample("bad", map[string]int{"count": -1}, nil)

	if cursor != nil {
		t.Fatal("Expected a nil cursor.")
	}

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected 404 with bad collection name.",
	)
}

func TestByExampleGetsNothing(t *testing.T) {
	byExampleDataInit()
	var se = getSE("_system")

	cursor, err := se.ByExample("test", map[string]int{"count": -1}, nil)

	if err != nil {
		t.Fatal("Uexpected error when searching by example: ", err)
	}

	if cursor.Id() != "" {
		t.Fatalf(
			"Expected cursor to NOT have an Id (%s)",
			cursor.Id(),
		)
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

	if cursor.Id() != "" {
		t.Fatalf(
			"Expected cursor to NOT have an Id (%s)",
			cursor.Id(),
		)
	}

	if !cursor.HasMore() {
		t.Fatal("Expected : cursor.HasMore() == true")
	}

	if cursor.Count() != 1 {
		t.Fatalf(
			"Actual cursor.Count() == %d, Expected : cursor.Count() == 1",
			cursor.Count(),
		)
	}

	var fetcher countDocument
	err = cursor.Next(&fetcher)

	if err != nil {
		t.Fatal("Unexpected error when reading Next:", err)
	}

	if fetcher.Count != 1 {
		t.Fatal("Fetched the wrong document: ", fetcher)
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

func TestByExampleFetchesMultiple(t *testing.T) {
	var se = getSE("_system")

	cursor, err := se.ByExample("test", map[string]string{"letter": "a"}, nil)

	if err != nil {
		t.Fatal("Uexpected error when searching by example: ", err)
	}

	if cursor.Id() != "" {
		t.Fatalf(
			"Expected cursor to NOT have an Id (%s)",
			cursor.Id(),
		)
	}

	if !cursor.HasMore() {
		t.Fatal("Expected : cursor.HasMore() == true")
	}

	if cursor.Count() != 5 {
		t.Fatalf(
			"Actual cursor.Count() == %d, Expected : cursor.Count() == 5",
			cursor.Count(),
		)
	}

	var fetcher countDocument
	var seenCounts = map[int]struct{}{}

	for cursor.HasMore() {
		err = cursor.Next(&fetcher)
		if err != nil {
			t.Fatal("Unexpected error when reading Next:", err)
		}
		if fetcher.Count > 4 {
			t.Fatal("Fetched the wrong document.")
		}
		if _, ok := seenCounts[fetcher.Count]; ok {
			t.Fatal(
				"Fetched the same document twice for same reason: %s",
				fetcher,
			)
		}

		seenCounts[fetcher.Count] = struct{}{}

		if fetcher.Letter != "a" {
			t.Fatal("Expected letter to be filled.")
		}
	}

	if cursor.HasMore() {
		t.Fatal("Expected : cursor.HasMore() == false")
	}

	err = cursor.Next(&fetcher)

	if err != EOC {
		t.Fatal("Expected end of cursor error when trying to read next.")
	}
}

func TestSimpleByExampleBadLimit(t *testing.T) {
	var se = getSE("_system")

	cursor, err := se.ByExample(
		"test",
		map[string]string{"letter": "b"},
		&ByExampleOptions{
			Limit:     -1,
			BatchSize: 1,
		},
	)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected an error with negative numbers for Limit",
	)

	if cursor != nil {
		t.Fatal("Expected cursor to be nil.")
	}
}

func TestSimpleByExampleBadBatchSize(t *testing.T) {
	var se = getSE("_system")

	cursor, err := se.ByExample(
		"test",
		map[string]string{"letter": "b"},
		&ByExampleOptions{
			Limit:     2,
			BatchSize: -1,
		},
	)

	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	if cursor == nil {
		t.Fatal("Expected a cursor even with negative batch size.")
	}

	if cursor.Count() != 2 {
		t.Fatalf(
			"Actual cursor.Count() == %d, Expected cursor.Count() == 2",
			cursor.Count(),
		)
	}

}

func TestSimpleByExampleLimitAndBatchSize(t *testing.T) {
	var se = getSE("_system")

	cursor, err := se.ByExample(
		"test",
		map[string]string{"letter": "b"},
		&ByExampleOptions{
			Limit:     2,
			BatchSize: 1,
		},
	)

	if cursor.Id() == "" {
		t.Fatal(
			"Expected cursor to have an Id but it didn't.",
		)
	}

	if err != nil {
		t.Fatal("Uexpected error when searching by example: ", err)
	}

	if !cursor.HasMore() {
		t.Fatal("Expected : cursor.HasMore() == true")
	}

	if cursor.Count() != 2 {
		t.Fatalf(
			"Actual cursor.Count() == %d, Expected : cursor.Count() == 2",
			cursor.Count(),
		)
	}

	var fetcher countDocument
	var seenCounts = map[int]struct{}{}

	for cursor.HasMore() {
		err = cursor.Next(&fetcher)

		if err != nil {
			t.Fatal("Unexpected error when reading Next:", err)
		}

		seenCounts[fetcher.Count] = struct{}{}

		if fetcher.Letter != "b" {
			t.Fatal("Expected letter to be filled.")
		}
	}

	if len(seenCounts) != 2 {
		t.Fatal("Cursor did not have a limit of 2 documents.")
	}

	if cursor.HasMore() {
		t.Fatal("Expected : cursor.HasMore() == false")
	}

	err = cursor.Next(&fetcher)

	if err != EOC {
		t.Fatal("Expected end of cursor error when trying to read next.")
	}

}

func TestFinalByExampleTest(t *testing.T) {
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
		if i >= 5 {
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
