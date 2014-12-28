package arango

import (
	"testing"
)

func TestSearchByExample(t *testing.T) {
	setup()
	defer teardown()

	db := db

	type internal struct {
		Field string `json:"field"`
	}

	type basic struct {
		DocumentImplementation
		Field string   `json:"field"`
		Obj   internal `json:"obj"`
	}

	type basicEdge struct {
		EdgeImplementation
		Field string `json:"field"`
	}

	d, err := db.CreateDocumentCollection("simple_docs")
	//e, err := db.CreateEdgeCollection( "simple_edges" )

	var doc1 = basic{Field: "hi", Obj: internal{Field: "hi"}}
	d.Save(&doc1)

	var doc2 = basic{Field: "bye", Obj: internal{Field: "bye"}}
	d.Save(&doc2)

	//Should return everything
	cur, err := d.ByExample(&struct{}{})

	if err != nil {
		t.Fatal("Cursor not returned but we exected everything to be fine.")
	}

	if cur.Error() {
		t.Fatal("Cursor fetched has an internal error that is not expected.", cur)
	}

	if cur.Count() != 2 {
		t.Fatalf("Cursor count not what we expected: %d", cur)
	}

	if !cur.HasMore() {
		t.Fatal("Expected cursor to have next.", cur)
	}

	var fetchDoc = basic{}

	i := 0
	for cur.HasMore() {
		i++
		err = cur.Next(&fetchDoc)
		if err != nil {
			t.Fatal("Error while fetching the next document", err)
		}
		switch fetchDoc.Id() {
		case doc1.Id():
			if fetchDoc.Rev() != doc1.Rev() ||
				fetchDoc.Field != doc1.Field ||
				fetchDoc.Obj.Field != doc1.Obj.Field {
				t.Fatalf("The fetched doc and the original are not equal.")
			}
		case doc2.Id():
			if fetchDoc.Rev() != doc2.Rev() ||
				fetchDoc.Field != doc2.Field ||
				fetchDoc.Obj.Field != doc2.Obj.Field {
				t.Fatalf("The fetched doc and the original are not equal.")
			}
		default:
			t.Fatal("Found a document that was not expected. No idea what to do.")
		}
		fetchDoc = basic{}
		if i > 2 {
			t.Fatalf("Got stuck in a loop but should've only gotten two documents: %+v", cur)
		}
	}

	err = cur.Next(&fetchDoc)

	if err == nil {
		t.Fatal("Expected error when calling a cursor with nothing left.")
	}

	cur, err = d.ByExampleQuery(&ByExampleQuery{
		Example:   &struct{}{},
		BatchSize: 1,
	})

	fetchDoc = basic{}
	i = 0
	for cur.HasMore() {
		i++
		err = cur.Next(&fetchDoc)
		if err != nil {
			t.Fatal("Error while fetching the next document", err)
		}
		switch fetchDoc.Id() {
		case doc1.Id():
			if fetchDoc.Rev() != doc1.Rev() ||
				fetchDoc.Field != doc1.Field ||
				fetchDoc.Obj.Field != doc1.Obj.Field {
				t.Fatalf("The fetched doc and the original are not equal.")
			}
		case doc2.Id():
			if fetchDoc.Rev() != doc2.Rev() ||
				fetchDoc.Field != doc2.Field ||
				fetchDoc.Obj.Field != doc2.Obj.Field {
				t.Fatalf("The fetched doc and the original are not equal.")
			}
		default:
			t.Fatal("Found a document that was not expected. No idea what to do.")
		}
		fetchDoc = basic{}
		if i > 2 {
			t.Fatalf("Got stuck in a loop but should've only gotten two documents: %+v", cur)
		}
	}

	err = cur.Next(&fetchDoc)

	if err == nil {
		t.Fatal("Expected error when calling a cursor with nothing left.")
	}

}
