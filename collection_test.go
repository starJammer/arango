package arango

import (
	"net/url"
	"testing"
)

func TestGetCollection(t *testing.T) {

	var db Database = getDatabase("_system")
	var ce CollectionEndpoint = db.CollectionEndpoint()

	if ce.Database() == nil {
		t.Fatal("Expected a link back to collections database.")
	}

}

func getCE(database string) CollectionEndpoint {
	var db Database = getDatabase(database)
	return db.CollectionEndpoint()
}

func TestGetCollections(t *testing.T) {
	var ce = getCE("_system")

	collections, err := ce.GetCollections(false)

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

func TestPostGetCollection(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")
	var ce = db.CollectionEndpoint()

	err := ce.PostCollection("", nil)

	if err == nil {
		t.Fatal("Expected error whet creating collection with no name or options.")
	}

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"

	err = ce.PostCollection(opts.Name, nil)

	if err != nil {
		t.Fatal("Unexpected error when creating collection: ", err)
	}

	colls, err := ce.GetCollections(true)

	if found := colls.Find(opts.Name); found == nil || found.Name() != opts.Name {
		t.Fatal("Could not find newly created connection.")
	}

	descriptor, err := ce.Get(opts.Name)

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

	err = ce.Delete(opts.Name)

	if err != nil {
		t.Fatal("Unexpected error when deleting collection: ", err)
	}

	opts.Type = EDGE_COLLECTION
	opts.Name = "test"

	err = ce.PostCollection(opts.Name, opts)

	descriptor, err = ce.Get(opts.Name)

	if err != nil {
		t.Fatal("Unexpected result from CollectionEndpoint.Get", err)
	}

	if descriptor.Type() != EDGE_COLLECTION {
		t.Fatal("Expected collection to be of type EDGE: ", descriptor.Type())
	}

	err = ce.Delete(opts.Name)
	if err != nil {
		t.Fatal("Unexpected error when deleting collection: ", err)
	}
}

func TestGetCollectionProperties(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")
	var ce = db.CollectionEndpoint()

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.WaitForSync = true
	opts.DoCompact = false

	err := ce.PostCollection(opts.Name, opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer ce.Delete(opts.Name)

	descriptor, err := ce.GetProperties(opts.Name)

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
	var ce = db.CollectionEndpoint()

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"

	err := ce.PostCollection(opts.Name, opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer ce.Delete(opts.Name)

	descriptor, err := ce.GetCount(opts.Name)

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
	var ce = db.CollectionEndpoint()

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"

	err := ce.PostCollection(opts.Name, opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer ce.Delete(opts.Name)

	descriptor, err := ce.GetFigures(opts.Name)

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
	var ce = db.CollectionEndpoint()

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"

	err := ce.PostCollection(opts.Name, opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer ce.Delete(opts.Name)

	descriptor, err := ce.GetRevision(opts.Name)

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
	var ce = db.CollectionEndpoint()

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"

	err := ce.PostCollection(opts.Name, opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer ce.Delete(opts.Name)

	descriptor, err := ce.GetChecksum(opts.Name, false, false)

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
	var ce = db.CollectionEndpoint()

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"

	err := ce.PostCollection(opts.Name, opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer ce.Delete(opts.Name)

	descriptor, err := ce.PutLoad(opts.Name, false)

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
	var ce = db.CollectionEndpoint()

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"

	err := ce.PostCollection(opts.Name, opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer ce.Delete(opts.Name)

	descriptor, err := ce.PutUnload(opts.Name)

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
	var ce = db.CollectionEndpoint()

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"

	err := ce.PostCollection(opts.Name, opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer ce.Delete(opts.Name)
	_, err = ce.PutTruncate(opts.Name)

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

}

func TestPutProperties(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")
	var ce = db.CollectionEndpoint()

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"
	opts.WaitForSync = true

	err := ce.PostCollection(opts.Name, opts)

	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}
	defer ce.Delete(opts.Name)

	descriptor, err := ce.GetProperties(opts.Name)
	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	if descriptor.WaitForSync() != true {
		t.Fatal("Expected waitforsync to be true upon creation.")
	}

	descriptor, err = ce.PutProperties(opts.Name, &CollectionPropertyChange{WaitForSync: false})

	if err != nil {
		t.Fatal("Unexpected error when getting collection descriptor: ", err)
	}

	descriptor, err = ce.PutProperties(opts.Name, &CollectionPropertyChange{WaitForSync: false})

	if descriptor.WaitForSync() != false {
		t.Fatal("Expected waitforSync to be false now.")
	}
}

func TestPutRename(t *testing.T) {
	u, _ := url.Parse("http://root@localhost:8529")
	c, _ := NewConnection(u)
	var db Database = c.Database("_system")
	var ce = db.CollectionEndpoint()
	var newName = "newtestname"

	opts := DefaultPostCollectionOptions()
	opts.Name = "test"

	err := ce.PostCollection(opts.Name, opts)
	if err != nil {
		t.Fatal("Error creating collection: ", err)
	}

	descriptor, err := ce.PutRename(opts.Name, newName)

	if err != nil {
		t.Fatal("Error during rename: ", err)
	}

	if descriptor.Name() != newName {
		t.Fatal("Collection rename failed: ", descriptor.Name())
	}

	_, err = ce.GetProperties(opts.Name)

	if err == nil {
		t.Fatal("Expected an error when getting properties of old collection.")
	}

	err = ce.Delete(newName)

	if err != nil {
		t.Fatal("Error deleting newly renamed collection: ", err)
	}
}
