package arango

import(
    "testing"
)

func TestSearchByExample( t *testing.T ) {
	setup()
	defer teardown()

	db := db

    type internal struct {
        Field string `json:"field"`
    }

	type basic struct {
		DocumentImplementation
        Field string `json:"field"`
        Obj internal `json:"obj"`
	}

	type basicEdge struct {
		EdgeImplementation
		Field string `json:"field"`
	}

    d, err := db.CreateDocumentCollection( "simple_docs" )
    //e, err := db.CreateEdgeCollection( "simple_edges" )

    var doc1 = basic{ Field : "hi", Obj : internal{ Field : "hi" } }
    d.Save( &doc1 )

    var doc2 = basic{ Field : "bye", Obj : internal{ Field : "bye" } }
    d.Save( &doc2 )

    //Should return everything
    cur, err := d.ByExample( &struct{}{} )

    if err != nil {
        t.Fatal( "Cursor not returned but we exected everything to be fine.")
    }

    if cur.Error() {
        t.Fatal( "Cursor fetched has an internal error that is not expected.", cur )
    }

    if cur.Count() != 2 {
        t.Fatalf( "Cursor count not what we expected: %d", cur )
    }

    if !cur.HasMore() {
        t.Fatal( "Expected cursor to have next.", cur )
    }

}
