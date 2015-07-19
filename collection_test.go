package arango

import (
	"net/url"
	"testing"
)

func TestGetCollection(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")

	opts := DefaultCollectionOptions()
	opts.Name = "test"

	err := db.PostCollection(opts)

	if err != nil {
		t.Fatal("Error creating collection")
	}
	defer db.Collection(opts.Name).Delete()

	col := db.Collection(opts.Name)

	descriptor, err := col.Get()

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	if descriptor.Id() == "" {
		t.Fatal("No id present in newly created collection.")
	}

	if descriptor.Name() != opts.Name {
		t.Fatalf("Unexpected collection name - Expected(%s) Actual(%s)", opts.Name, descriptor.Name())
	}

	if descriptor.Status() != LOADED_STATUS {
		t.Fatalf("Unexpected collection status - Expected(%d) Actual(%d)", LOADED_STATUS, descriptor.Status())
	}

	if descriptor.Type() != DOCUMENT_COLLECTION {
		t.Fatalf("Unexpected collection type - Expected(%d) Actual(%d)", DOCUMENT_COLLECTION, descriptor.Type())
	}

	if descriptor.IsSystem() != false {
		t.Fatalf("Unexpected IsSystem value - Expected(%t) Actual(%t)", false, descriptor.IsSystem())
	}

}

func TestGetCollectionProperties(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")

	opts := DefaultCollectionOptions()
	opts.Name = "test"
	opts.WaitForSync = true
	opts.DoCompact = false

	err := db.PostCollection(opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer db.Collection(opts.Name).Delete()

	col := db.Collection(opts.Name)

	descriptor, err := col.GetProperties()

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	if descriptor.WaitForSync() != true {
		t.Fatal("Expected waitForSync to be true.")
	}

	if descriptor.DoCompact() != false {
		t.Fatal("Expected doCompact to be false.")
	}

	if descriptor.JournalSize() == 0 {
		t.Fatal("Expected positive value for journal size.")
	}

	if descriptor.KeyOptions() == nil {
		t.Fatal("Expected some key options.")

		k := descriptor.KeyOptions()
		if k.Type() == "" {
			t.Fatal("Expected a value for keyoptions type")
		}
	}

	if descriptor.IsVolatile() != false {
		t.Fatal("Expected isVolatile to be false.")
	}

	if descriptor.NumberOfShards() != 0 {
		t.Fatal("Expected numberOfShards to be 0: ", descriptor.NumberOfShards())
	}

	if len(descriptor.ShardKeys()) != 0 {
		t.Fatal("Expected ShardKeys to have a length of zero.: ", descriptor.ShardKeys())
	}

}

func TestGetCollectionCount(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")

	opts := DefaultCollectionOptions()
	opts.Name = "test"

	err := db.PostCollection(opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer db.Collection(opts.Name).Delete()
	col := db.Collection(opts.Name)

	descriptor, err := col.GetCount()

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	if descriptor.Count() != 0 {
		t.Fatal("Expected to get a count of 0.")
	}

}

func TestGetCollectionFigures(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")

	opts := DefaultCollectionOptions()
	opts.Name = "test"

	err := db.PostCollection(opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer db.Collection(opts.Name).Delete()
	col := db.Collection(opts.Name)

	descriptor, err := col.GetFigures()

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	if descriptor.Figures() == nil {
		t.Fatal("Expected figures to be non-nil.")
	}

}

func TestGetCollectionRevision(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")

	opts := DefaultCollectionOptions()
	opts.Name = "test"

	err := db.PostCollection(opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer db.Collection(opts.Name).Delete()
	col := db.Collection(opts.Name)

	descriptor, err := col.GetRevision()

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	if descriptor.Revision() == "" {
		t.Fatal("Expected revision to be non-blank.")
	}

}

func TestGetCollectionChecksum(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")

	opts := DefaultCollectionOptions()
	opts.Name = "test"

	err := db.PostCollection(opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer db.Collection(opts.Name).Delete()
	col := db.Collection(opts.Name)

	descriptor, err := col.GetChecksum(false, false)

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	//checksum for new collections is 0.
	//Test checksum for collections with docs within
	//later
	if descriptor.Checksum() != 0 {
		t.Fatal("Expected checksum to be non-zero.")
	}

}
