package arango

import (
    "testing"
)

func TestDatabaseMethods( t *testing.T ) {
    db, err := Conn( "http://root@localhost:8529")
    if err != nil {
        t.Fatal( err )
    }

    err = db.DropDatabase( "testing" )

    if err != nil {
        e, _ := err.(ArangoError)
        if e.Code != 404 && e.Code == 200 {
            t.Log( "Database exists but was expecting it not to.", err )
        }
    }

    err = db.CreateDatabase( "testing", nil, nil )

    if err != nil { 
        t.Fatal( err )
    }
}
