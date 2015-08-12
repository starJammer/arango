package arango

import (
	"net/http"
	"testing"
)

func verifyError(err error, t *testing.T, code int, message string) {
	if err == nil {
		t.Fatal(message)
	}
	ae, ok := err.(ArangoError)
	if !ok {
		t.Fatalf("Expected an ArangoError to be returned. %#v", err)
	}
	if !ae.IsError() {
		t.Fatalf("Actual ae.IsError() == %t, Expected true. ArangoError = %s, Message = %s", ae.IsError(), ae, message)
	}

	if ae.Code() != code {
		t.Fatalf("Actual ae.Code() == %d, Expected %d. ArangoError = %s, Message = %s", ae.Code(), code, ae, message)
	}
}

func TestGetDatabase(t *testing.T) {
	c := setupConnection()

	var db Database = c.Database("_system")

	if db.Name() != "_system" {
		t.Fatal("Database name incorrect: ", db.Name())
	}

}

func getDatabase(name string) Database {
	c := setupConnection()
	return c.Database(name)
}

func TestGetDatabases(t *testing.T) {
	db := getDatabase("_system")

	dbs, err := db.Get()

	if err != nil {
		t.Fatal("Could not fetch all the databases: ", err)
	}

	if len(dbs) < 1 {
		t.Fatal("Expected to at least have the _system database but don't.")
	}

	userDbs, err := db.GetUser()

	if err != nil {
		t.Fatal("Could not fetch user databases: ", err)
	}

	if len(userDbs) != len(dbs) {
		t.Fatalf("User dbs not the same as regular dbs - user: %s, regular: %s", userDbs, dbs)
	}
}

func TestGetCurrent(t *testing.T) {
	db := getDatabase("_system")

	current, err := db.GetCurrent()

	if err != nil {
		t.Fatal("Could not fetch current db info:", err)
	}

	if current.Name() != "_system" {
		t.Fatal("Unexpected name: ", current.Name())
	}

	if current.Id() == "" {
		t.Fatal("Unexpected id: ", current.Id())
	}

	if current.Path() == "" {
		t.Fatal("Unexpected path: ", current.Path())
	}

	if !current.IsSystem() {
		t.Fatalf("Unexpected isSystem: %t", current.IsSystem())
	}
}

func TestGetCurrentForBadDatabase(t *testing.T) {

	var bad Database = getDatabase("baddatabase")
	current, err := bad.GetCurrent()
	if current != nil {
		t.Fatalf("Expected current to be nil: %#v", current)
	}
	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected an error with GetCurrent on non-existent database.",
	)

}

func TestPostWithBlankName(t *testing.T) {
	var db Database = getDatabase("_system")
	var err error

	err = db.Post("", nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error creating a database with no name.",
	)
}

func TestPostWithExistingDb(t *testing.T) {

	var db Database = getDatabase("_system")
	var err error

	err = db.Post("_system", nil)

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when creating a database that exists already.",
	)
}

func TestPostGoodNameThenDelete(t *testing.T) {

	var db Database = getDatabase("_system")
	var err error

	err = db.Post("test", nil)

	if err != nil {
		t.Fatal("Unexpected error when creating new database: ", err)
	}

	err = db.Delete("test")

	if err != nil {
		t.Fatal("Unexpected error when deleting a database: ", err)
	}
}

func TestPostGetCurrent(t *testing.T) {

	var system Database = getDatabase("_system")

	system.Post("test", nil)
	defer system.Delete("test")

	var test Database = getDatabase("test")
	current, _ := test.GetCurrent()

	if current.Name() != "test" {
		t.Fatal("Unexpected name: ", current.Name())
	}

	if current.Id() == "" {
		t.Fatal("Unexpected id: ", current.Id())
	}

	if current.Path() == "" {
		t.Fatal("Unexpected path: ", current.Path())
	}

	if current.IsSystem() {
		t.Fatalf("Unexpected isSystem: %t", current.IsSystem())
	}

}

func TestPostDuplicateDatabaseName(t *testing.T) {

	var db Database = getDatabase("_system")
	var err error

	_ = db.Post("test", nil)
	defer db.Delete("test")

	err = db.Post("test", nil)
	verifyError(
		err,
		t,
		http.StatusConflict,
		"Expected error when creating a database with the same name.",
	)
}

func TestDeleteBlankName(t *testing.T) {

	var db Database = getDatabase("_system")
	var err error

	err = db.Delete("")
	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when deleting database with a blank name.",
	)
}

func TestDeleteBadName(t *testing.T) {

	var db Database = getDatabase("_system")
	var err error

	err = db.Delete("badname")
	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when deleting database with a bad name.",
	)

}
