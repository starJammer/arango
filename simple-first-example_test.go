package arango

import (
	"net/http"
	"testing"
)

func TestFirstExampleBlankCollection(t *testing.T) {
	var se = getSE("_system")

	var fetcher *countDocument

	err := se.FirstExample("", map[string]int{"count": -1}, fetcher)

	if fetcher != nil {
		t.Fatal("Expected a nil fetcher.")
	}

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected 404 with blank collection name.",
	)
}

func TestFirstExampleBadCollection(t *testing.T) {
	var se = getSE("_system")

	var fetcher *countDocument

	err := se.FirstExample("bad", map[string]int{"count": -1}, fetcher)

	if fetcher != nil {
		t.Fatal("Expected a nil fetcher.")
	}

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected 404 with bad collection name.",
	)
}

func TestFirstExampleGetsNothing(t *testing.T) {
	byExampleDataInit()
	var se = getSE("_system")

	var fetcher *countDocument

	err := se.FirstExample("test", map[string]int{"count": -1}, fetcher)

	if fetcher != nil {
		t.Fatal("Expected a nil fetcher.")
	}

	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected no matches for the search.",
	)

}

func TestFirstExample(t *testing.T) {
	var se = getSE("_system")

	var fetcher *countDocument
	err := se.FirstExample("test", map[string]int{"count": 1}, &fetcher)

	if err != nil {
		t.Fatal("Expected no error for FirstExample by count == 1.", err)
	}

	if fetcher.Count != 1 {
		t.Fatal(
			"Actual fetcher.Count == %d, Expected fetcher.Count == 1",
			fetcher.Count,
		)
	}
}

func TestFirstExampleAmbiguousSearch(t *testing.T) {
	var se = getSE("_system")

	var fetcher *countDocument
	err := se.FirstExample("test", map[string]string{"letter": "a"}, &fetcher)

	if err != nil {
		t.Fatal("Expected no error for FirstExample by letter == a: ", err)
	}

	if fetcher == nil {
		t.Fatal("Unexpected nil value for fetcher.")
	}

	if fetcher.Letter != "a" {
		t.Fatal(
			"Actual fetcher.Letter == %s, Expected fetcher.Letter == \"a\"",
			fetcher.Letter,
		)
	}

}

func TestFinalFirstExampleTest(t *testing.T) {
	destroyByExampleDataInit()
}
