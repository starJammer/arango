package arango

import (
	"net/url"
	"testing"
)

func TestGetCollections(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")
	var collEnd = db.CollectionEndpoint()

	collections, err := collEnd.GetCollections(false)

	if err != nil {
		t.Fatal("Unexpected error during GetCollections: ", err)
	}

	if len(collections) < 1 {
		t.Fatal("Expected at least one collection.")
	}

	coll := collections[0]

	if coll.Id() == "" {
		t.Fatal("Expected an id value for the collection.")
	}

	if coll.Name() == "" {
		t.Fatal("Expected a name value for the collection.")
	}

	if coll.Status() == 0 {
		t.Fatal("Expected a CollectionStatus value for the collection.")
	}

	if coll.Type() == 0 {
		t.Fatal("Expected a CollectionType value for the collection.")
	}

}

func TestPostCollection(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")
	var collEnd = db.CollectionEndpoint()

	err := collEnd.PostCollection(nil)

	if err == nil {
		t.Fatal("Expected error whet creating collection with nil options.")
	}

	opts := DefaultCollectionOptions()
	opts.Name = "test"

	err = collEnd.PostCollection(opts)

	if err != nil {
		t.Fatal("Got an error when creating collection: ", err)
	}

	colls, err := collEnd.GetCollections(true)

	if found := colls.Find(opts.Name); found == nil || found.Name() != opts.Name {
		t.Fatal("Could not find newly created connection.")
	}

	err = collEnd.Delete(opts.Name)

	if err != nil {
		t.Fatal("Unexpected error when deleting collection: ", err)
	}
}

func TestGetCollection(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")
	var collEnd = db.CollectionEndpoint()

	opts := DefaultCollectionOptions()
	opts.Name = "test"

	err := collEnd.PostCollection(opts)

	if err != nil {
		t.Fatal("Error creating collection")
	}
	defer collEnd.Delete(opts.Name)

	descriptor, err := collEnd.Get(opts.Name)

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
	var collEnd = db.CollectionEndpoint()

	opts := DefaultCollectionOptions()
	opts.Name = "test"
	opts.WaitForSync = true
	opts.DoCompact = false

	err := collEnd.PostCollection(opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer collEnd.Delete(opts.Name)

	descriptor, err := collEnd.GetProperties(opts.Name)

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
	var collEnd = db.CollectionEndpoint()

	opts := DefaultCollectionOptions()
	opts.Name = "test"

	err := collEnd.PostCollection(opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer collEnd.Delete(opts.Name)

	descriptor, err := collEnd.GetCount(opts.Name)

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
	var collEnd = db.CollectionEndpoint()

	opts := DefaultCollectionOptions()
	opts.Name = "test"

	err := collEnd.PostCollection(opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer collEnd.Delete(opts.Name)

	descriptor, err := collEnd.GetFigures(opts.Name)

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
	var collEnd = db.CollectionEndpoint()

	opts := DefaultCollectionOptions()
	opts.Name = "test"

	err := collEnd.PostCollection(opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer collEnd.Delete(opts.Name)

	descriptor, err := collEnd.GetRevision(opts.Name)

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
	var collEnd = db.CollectionEndpoint()

	opts := DefaultCollectionOptions()
	opts.Name = "test"

	err := collEnd.PostCollection(opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer collEnd.Delete(opts.Name)

	descriptor, err := collEnd.GetChecksum(opts.Name, false, false)

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

func TestPutLoad(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")
	var collEnd = db.CollectionEndpoint()

	opts := DefaultCollectionOptions()
	opts.Name = "test"

	err := collEnd.PostCollection(opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer collEnd.Delete(opts.Name)

	descriptor, err := collEnd.PutLoad(opts.Name, false)

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	if descriptor.Status() != LOADED_STATUS {
		t.Fatal("Expected collection to be in loaded state.")
	}

}

func TestPutUnload(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")
	var collEnd = db.CollectionEndpoint()

	opts := DefaultCollectionOptions()
	opts.Name = "test"

	err := collEnd.PostCollection(opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer collEnd.Delete(opts.Name)

	descriptor, err := collEnd.PutUnload(opts.Name)

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	if descriptor.Status() == LOADED_STATUS {
		t.Fatal("Expected collection to not be in loaded state: ", descriptor.Status())
	}
}

func TestPutTruncate(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")
	var collEnd = db.CollectionEndpoint()

	opts := DefaultCollectionOptions()
	opts.Name = "test"

	err := collEnd.PostCollection(opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer collEnd.Delete(opts.Name)
	_, err = collEnd.PutTruncate(opts.Name)

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

}

func TestPutProperties(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")
	var collEnd = db.CollectionEndpoint()

	opts := DefaultCollectionOptions()
	opts.Name = "test"
	opts.WaitForSync = true

	err := collEnd.PostCollection(opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer collEnd.Delete(opts.Name)

	descriptor, err := collEnd.GetProperties(opts.Name)
	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	if descriptor.WaitForSync() != true {
		t.Fatal("Expected waitforsync to be true upon creation.")
	}

	descriptor, err = collEnd.PutProperties(opts.Name, &CollectionPropertyChange{WaitForSync: false})

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	descriptor, err = collEnd.PutProperties(opts.Name, &CollectionPropertyChange{WaitForSync: false})

	if descriptor.WaitForSync() != false {
		t.Fatal("Expected waitforSync to be false now.")
	}
}

func TestPutRename(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")
	var collEnd = db.CollectionEndpoint()
	var newName = "newtestname"

	opts := DefaultCollectionOptions()
	opts.Name = "test"

	err := collEnd.PostCollection(opts)
	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}

	descriptor, err := collEnd.PutRename(opts.Name, newName)

	if err != nil {
		t.Fatal("Error during rename: ", err)
	}

	if descriptor.Name() != newName {
		t.Fatal("Collection rename failed: ", descriptor.Name())
	}

	_, err = collEnd.GetProperties(opts.Name)

	if err == nil {
		t.Fatal("Expected an error when getting properties of old collection.")
	}

	err = collEnd.Delete(newName)

	if err != nil {
		t.Fatal("Error deleting newly renamed collection: ", err)
	}
}
