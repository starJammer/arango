package arango

import (
    "testing"
)

func TestDatabaseMethods( t *testing.T ) {
    _, _ := Conn( "http://root@localhost:8529")
}
