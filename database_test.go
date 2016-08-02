package arango

import (
	"net/http"
	"testing"
)

func TestGetDatabase(t *testing.T) {
	c := setupConnection()

	var db = c.Database("_system")

	if db.Name() != "_system" {
		t.Fatal("Database name incorrect: ", db.Name())
	}

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

	if current.Name != "_system" {
		t.Fatal("Unexpected name: ", current.Name)
	}

	if current.Id == "" {
		t.Fatal("Unexpected id: ", current.Id)
	}

	if current.Path == "" {
		t.Fatal("Unexpected path: ", current.Path)
	}

	if !current.IsSystem {
		t.Fatalf("Unexpected isSystem: %t", current.IsSystem)
	}
}

func TestGetCurrentForBadDatabase(t *testing.T) {

	var bad = getDatabase("baddatabase")
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
	var db = getDatabase("_system")
	var err error

	err = db.Post(&PostDatabaseOptions{
		Name: "",
	})

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error creating a database with no name.",
	)
}

func TestPostWithExistingDb(t *testing.T) {

	var db = getDatabase("_system")
	var err error

	err = db.Post(&PostDatabaseOptions{
		Name: "_system",
	})

	verifyError(
		err,
		t,
		http.StatusBadRequest,
		"Expected error when creating a database that exists already.",
	)
}

func TestCantPostOutsideSystem(t *testing.T) {

	var db = getDatabase("_system")
	var err error

	db.Post(&PostDatabaseOptions{
		Name: "test",
	})
	defer db.Delete("test")

	testDb := getDatabase("test")

	err = testDb.Post(&PostDatabaseOptions{
		Name: "fail",
	})

	verifyError(
		err,
		t,
		http.StatusForbidden,
		"Expected to receive error when posting database outside of _system",
	)
}

func TestPostGoodNameThenDelete(t *testing.T) {

	var db = getDatabase("_system")
	var err error

	err = db.Post(&PostDatabaseOptions{
		Name: "test",
	})

	if err != nil {
		t.Fatal("Unexpected error when creating new database: ", err)
	}

	err = db.Delete("test")

	if err != nil {
		t.Fatal("Unexpected error when deleting a database: ", err)
	}
}

func TestPostGetCurrent(t *testing.T) {

	var system = getDatabase("_system")

	system.Post(&PostDatabaseOptions{
		Name: "test",
	})
	defer system.Delete("test")

	var test = getDatabase("test")
	current, _ := test.GetCurrent()

	if current.Name != "test" {
		t.Fatal("Unexpected name: ", current.Name)
	}

	if current.Id == "" {
		t.Fatal("Unexpected id: ", current.Id)
	}

	if current.Path == "" {
		t.Fatal("Unexpected path: ", current.Path)
	}

	if current.IsSystem {
		t.Fatalf("Unexpected isSystem: %t", current.IsSystem)
	}

}

func TestPostDuplicateDatabaseName(t *testing.T) {

	var db = getDatabase("_system")
	var err error
	var opts = &PostDatabaseOptions{
		Name: "test",
	}

	_ = db.Post(opts)
	defer db.Delete("test")

	err = db.Post(opts)
	verifyError(
		err,
		t,
		http.StatusConflict,
		"Expected error when creating a database with the same name.",
	)
}

func TestDeleteBlankName(t *testing.T) {

	var db = getDatabase("_system")
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

	var db = getDatabase("_system")
	var err error

	err = db.Delete("badname")
	verifyError(
		err,
		t,
		http.StatusNotFound,
		"Expected error when deleting database with a bad name.",
	)

}

func TestDeleteOutsideOfSystem(t *testing.T) {
	var db = getDatabase("_system")
	var err error

	db.Post(&PostDatabaseOptions{
		Name: "test",
	})
	defer db.Delete("test")

	testDb := getDatabase("test")

	err = testDb.Delete("test")

	verifyError(
		err,
		t,
		http.StatusForbidden,
		"Expected to receive error when posting database outside of _system",
	)
}
